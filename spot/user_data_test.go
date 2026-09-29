package spot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
)

// TestListenKeyKeepaliveCloseIgnoreBody checks that renewing and deleting the
// listenKey succeed on an HTTP 200 whatever the body, including a zero-byte
// one, and still surface API errors.
func TestListenKeyKeepaliveCloseIgnoreBody(t *testing.T) {
	for _, tt := range []struct {
		contentType string
		status      int
		body        string
		wantErr     bool
	}{
		{status: http.StatusOK},
		{contentType: "application/json", status: http.StatusOK},
		{contentType: "application/json", status: http.StatusOK, body: `{}`},
		{contentType: "application/json", status: http.StatusBadRequest, body: `{"code":-1125,"msg":"This listenKey does not exist."}`, wantErr: true},
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tt.contentType != "" {
				w.Header().Set("Content-Type", tt.contentType)
			}
			w.WriteHeader(tt.status)
			_, _ = w.Write([]byte(tt.body))
		}))
		c := NewSpotClient(client.WithBaseURL(srv.URL), client.WithAuth(testUserAddress, testSignerKeyHex))
		ctx := context.Background()
		errs := []error{
			c.NewRenewListenKeyService("key").Do(ctx),
			c.NewDeleteListenKeyService("key").Do(ctx),
		}
		srv.Close()
		for i, err := range errs {
			if (err != nil) != tt.wantErr {
				t.Errorf("content type %q, status %d, body %q, call %d: err = %v, wantErr %v", tt.contentType, tt.status, tt.body, i, err, tt.wantErr)
			}
		}
	}
}
