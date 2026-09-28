package request

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
)

// recordingLogger keeps the Errorf lines it is given.
type recordingLogger struct {
	mu     sync.Mutex
	errors []string
}

func (l *recordingLogger) Infof(string, ...any)  {}
func (l *recordingLogger) Warnf(string, ...any)  {}
func (l *recordingLogger) Debugf(string, ...any) {}
func (l *recordingLogger) Errorf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
}

// TestDoDecodeLogsFailures checks that DoDecode (and DoRaw) log transport,
// API and decode failures like Do, and nothing on success.
func TestDoDecodeLogsFailures(t *testing.T) {
	errDecode := errors.New("bad body")
	tests := []struct {
		name    string
		status  int
		decode  func([]byte) (string, error)
		wantErr bool
	}{
		{name: "ok", status: http.StatusOK, decode: func(b []byte) (string, error) { return string(b), nil }},
		{name: "decode error", status: http.StatusOK, decode: func([]byte) (string, error) { return "", errDecode }, wantErr: true},
		{name: "api error", status: http.StatusBadRequest, decode: func(b []byte) (string, error) { return string(b), nil }, wantErr: true},
	}
	for _, tt := range tests {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tt.status)
			_, _ = w.Write([]byte(`{"code":-1121,"msg":"Invalid symbol."}`))
		}))
		logger := &recordingLogger{}
		c := client.NewClient(client.ProductSpot, client.WithBaseURL(srv.URL), client.WithLogger(logger))
		got, err := DoDecode(Get(context.Background(), c, "/api/v3/x"), tt.decode)
		srv.Close()
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if tt.name == "ok" && got != `{"code":-1121,"msg":"Invalid symbol."}` {
			t.Errorf("%s: body = %q", tt.name, got)
		}
		if tt.name == "decode error" && !errors.Is(err, errDecode) {
			t.Errorf("%s: err = %v, want the decode error", tt.name, err)
		}
		if logged := len(logger.errors) == 1 && strings.Contains(logger.errors[0], "/api/v3/x"); logged != tt.wantErr {
			t.Errorf("%s: logged %q", tt.name, logger.errors)
		}
	}
}

// TestDoDecodesWhateverTheResponse checks that Do decodes the body whatever
// its Content-Type, never returns a nil *T with a nil error (a null body or a
// 204 gives the zero T), and still records the rate-limit headers and logs
// when decoding fails.
func TestDoDecodesWhateverTheResponse(t *testing.T) {
	type result struct {
		A int `json:"a"`
	}
	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
		want        int
		wantErr     bool
	}{
		{name: "json", status: http.StatusOK, contentType: "application/json", body: `{"a":1}`, want: 1},
		{name: "text/plain", status: http.StatusOK, contentType: "text/plain", body: `{"a":2}`, want: 2},
		{name: "no content type", status: http.StatusOK, body: `{"a":3}`, want: 3},
		{name: "null", status: http.StatusOK, contentType: "application/json", body: `null`},
		{name: "204", status: http.StatusNoContent},
		{name: "html", status: http.StatusOK, contentType: "text/html", body: `<html></html>`, wantErr: true},
		{name: "bad field", status: http.StatusOK, contentType: "application/json", body: `{"a":"x"}`, wantErr: true},
	}
	for _, tt := range tests {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tt.contentType != "" {
				w.Header().Set("Content-Type", tt.contentType)
			} else {
				w.Header()["Content-Type"] = nil
			}
			w.Header().Set("X-Mbx-Used-Weight-1m", "7")
			w.WriteHeader(tt.status)
			_, _ = w.Write([]byte(tt.body))
		}))
		logger := &recordingLogger{}
		c := client.NewClient(client.ProductSpot, client.WithBaseURL(srv.URL), client.WithLogger(logger))
		got, err := Do[result](Get(context.Background(), c, "/api/v3/x"))
		srv.Close()
		switch {
		case tt.wantErr && err == nil:
			t.Errorf("%s: got %+v, want error", tt.name, got)
		case !tt.wantErr && (err != nil || got == nil || got.A != tt.want):
			t.Errorf("%s: got %+v, %v, want {A:%d}", tt.name, got, err, tt.want)
		}
		if logged := len(logger.errors) == 1; logged != tt.wantErr {
			t.Errorf("%s: logged %q", tt.name, logger.errors)
		}
		if used := c.GetUsedWeight().Used1M; used != 7 {
			t.Errorf("%s: Used1M = %d, want 7", tt.name, used)
		}
	}
}
