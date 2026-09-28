package request

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/gorilla/websocket"
)

// serverMode is what wsServer does after sending its frames.
type serverMode int

const (
	serverReads    serverMode = iota // keeps reading, so it answers pings
	serverDrops                      // closes the connection
	serverStalls                     // goes silent: neither reads nor writes
	serverPings                      // pings the client but never reads
	serverTrickles                   // sends the frames 200ms apart, never reads
)

// wsServer sends frames to each connection, then behaves as mode says.
func wsServer(t *testing.T, frames []string, mode serverMode) *client.WebSocketClient {
	t.Helper()
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for _, f := range frames {
			if mode == serverTrickles {
				time.Sleep(200 * time.Millisecond)
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(f)); err != nil {
				return
			}
		}
		switch mode {
		case serverDrops:
			return
		case serverStalls, serverTrickles:
			time.Sleep(2 * time.Second)
			return
		case serverPings:
			for range 40 {
				time.Sleep(50 * time.Millisecond)
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
					return
				}
			}
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

// shortenKeepAlive makes subscriptions opened by the test ping every interval
// and give up after timeout of silence.
func shortenKeepAlive(t *testing.T, interval, timeout time.Duration) {
	t.Helper()
	i, to := keepAliveInterval, keepAliveTimeout
	keepAliveInterval, keepAliveTimeout = interval, timeout
	t.Cleanup(func() { keepAliveInterval, keepAliveTimeout = i, to })
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
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, serverReads), "/ws/x", rec.callback)
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
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, serverDrops), "/ws/x", rec.callback)
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
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, serverDrops), "/ws/x", rec.callback)
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

// TestSubscribeDetectsSilentPeer checks that a peer that stops sending
// anything, pongs included, ends the subscription with a timeout error.
func TestSubscribeDetectsSilentPeer(t *testing.T) {
	shortenKeepAlive(t, 50*time.Millisecond, 300*time.Millisecond)
	rec := &recorder{got: make(chan struct{})}
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, serverStalls), "/ws/x", rec.callback)
	if err != nil {
		t.Fatal(err)
	}
	waitClosed(t, stop, "stop")
	close(done)
	rec.mu.Lock()
	defer rec.mu.Unlock()
	var ne net.Error
	if len(rec.frames) != 2 || len(rec.errs) != 1 || !errors.As(rec.errs[0], &ne) || !ne.Timeout() {
		t.Errorf("frames %q, errors %v; want 2 frames and one timeout", rec.frames, rec.errs)
	}
}

// TestSubscribeKeepsQuietPeerAlive checks that a peer sending no data stays
// connected while it answers our pings or pings us itself.
func TestSubscribeKeepsQuietPeerAlive(t *testing.T) {
	for _, mode := range []serverMode{serverReads, serverPings} {
		shortenKeepAlive(t, 50*time.Millisecond, 300*time.Millisecond)
		rec := &recorder{got: make(chan struct{})}
		done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b"}, mode), "/ws/x", rec.callback)
		if err != nil {
			t.Fatal(err)
		}
		waitClosed(t, rec.got, "frames")
		select {
		case <-stop:
			t.Fatalf("mode %d: subscription ended on a quiet but live peer: %v", mode, rec.errs)
		case <-time.After(time.Second):
		}
		close(done)
		waitClosed(t, stop, "stop")
		rec.mu.Lock()
		if len(rec.errs) != 0 {
			t.Errorf("mode %d: errors %v, want none", mode, rec.errs)
		}
		rec.mu.Unlock()
	}
}

// TestSubscribeSlowCallback checks that time spent in callback does not count
// towards the read timeout: frames keep arriving faster than the timeout, so
// a callback that blocks longer than the timeout must not end the
// subscription. The server never reads, so no pong can extend the deadline.
func TestSubscribeSlowCallback(t *testing.T) {
	shortenKeepAlive(t, 50*time.Millisecond, 300*time.Millisecond)
	var mu sync.Mutex
	var frames []string
	var errs []error
	got := make(chan struct{})
	done, stop, err := SubscribeRaw(context.Background(), wsServer(t, []string{"a", "b", "c"}, serverTrickles), "/ws/x", func(msg []byte, err error) {
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		time.Sleep(500 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		if frames = append(frames, string(msg)); len(frames) == 3 {
			close(got)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-got:
	case <-stop:
	case <-time.After(5 * time.Second):
	}
	close(done)
	waitClosed(t, stop, "stop")
	mu.Lock()
	defer mu.Unlock()
	if len(frames) != 3 || len(errs) != 0 {
		t.Errorf("frames %q, errors %v; want 3 frames and no error", frames, errs)
	}
}
