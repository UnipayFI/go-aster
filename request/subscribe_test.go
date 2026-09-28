package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/gorilla/websocket"
)

// wsServer sends frames to each connection, then either waits for the client
// to go away or, with drop, closes the connection itself.
func wsServer(t *testing.T, frames []string, drop bool) *client.WebSocketClient {
	t.Helper()
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for _, f := range frames {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(f)); err != nil {
				return
			}
		}
		if drop {
			return
		}
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)
	return client.NewWebSocketClient(client.ProductSpot, client.WithWebSocketBaseURL("ws"+strings.TrimPrefix(srv.URL, "http")))
}

type recorder struct {
	mu     sync.Mutex
	frames []string
	errs   []error
	got    chan struct{}
}

func (r *recorder) callback(msg []byte, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err != nil {
		r.errs = append(r.errs, err)
		return
	}
	r.frames = append(r.frames, string(msg))
	if len(r.frames) == 2 {
		close(r.got)
	}
}

func waitClosed(t *testing.T, c <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-c:
	case <-time.After(5 * time.Second):
		t.Fatalf("%s was not closed", what)
	}
}

// TestSubscribeCallerClose checks that closing done ends the subscription
// without reporting the resulting read error, and closes stop.
func TestSubscribeCallerClose(t *testing.T) {
	rec := &recorder{got: make(chan struct{})}
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, false), "/ws/x", rec.callback)
	if err != nil {
		t.Fatal(err)
	}
	waitClosed(t, rec.got, "frames")
	close(done)
	waitClosed(t, stop, "stop")
	rec.mu.Lock()
	defer rec.mu.Unlock()
	if len(rec.frames) != 2 || len(rec.errs) != 0 {
		t.Errorf("frames %q, errors %v; want 2 frames and no error", rec.frames, rec.errs)
	}
}

// TestSubscribeServerDrop checks that a connection the server closes reports
// its error once and closes stop.
func TestSubscribeServerDrop(t *testing.T) {
	rec := &recorder{got: make(chan struct{})}
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, true), "/ws/x", rec.callback)
	if err != nil {
		t.Fatal(err)
	}
	waitClosed(t, stop, "stop")
	close(done)
	rec.mu.Lock()
	defer rec.mu.Unlock()
	if len(rec.frames) != 2 || len(rec.errs) != 1 {
		t.Errorf("frames %q, errors %v; want 2 frames and one error", rec.frames, rec.errs)
	}
}

// TestSubscribeSendOnDoneAfterDrop checks that sending on done, rather than
// closing it, does not block after the server has dropped the connection.
func TestSubscribeSendOnDoneAfterDrop(t *testing.T) {
	rec := &recorder{got: make(chan struct{})}
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, true), "/ws/x", rec.callback)
	if err != nil {
		t.Fatal(err)
	}
	waitClosed(t, stop, "stop")
	sent := make(chan struct{})
	go func() {
		done <- struct{}{}
		close(sent)
	}()
	waitClosed(t, sent, "send on done")
}
