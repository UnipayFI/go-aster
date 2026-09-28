package request

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

// keepAliveInterval is how often a subscription pings the server, and
// keepAliveTimeout how long it waits without any frame, ping or pong before
// giving the connection up. They are variables so tests can shorten them.
var (
	keepAliveInterval = common.DEFAULT_KEEP_ALIVE_INTERVAL
	keepAliveTimeout  = common.DEFAULT_KEEP_ALIVE_TIMEOUT
)

// controlWriteWait bounds the writes of ping and pong frames.
const controlWriteWait = 10 * time.Second

type WebSocketClient interface {
	GetHttpClient() *resty.Client
	GetLogger() log.Logger
	GetDialer() *websocket.Dialer
}

// Subscribe dials endpoint and decodes each frame into T for callback. Close
// (or send on) done to end the subscription; stop is closed once reading has
// ended, whether because done was closed or because a read failed, in which
// case callback gets the error first. Frames already read may still reach
// callback after done is closed, so wait on stop before releasing what
// callback uses.
func Subscribe[T any](ctx context.Context, client WebSocketClient, endpoint string, callback func(message *T, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, endpoint, func(message []byte, e error) {
		if e != nil {
			callback(nil, e)
			return
		}
		var msg T
		if err := client.GetHttpClient().JSONUnmarshal(message, &msg); err != nil {
			callback(nil, err)
			return
		}
		callback(&msg, nil)
	})
}

// SubscribeRaw delivers each WebSocket frame's raw bytes to the callback
// without JSON-decoding. Use this for streams that union multiple payload
// shapes under a single connection -- for example the spot/futures user-data
// stream, where the caller must look at the "e" field to decide which struct
// to decode into.
func SubscribeRaw(ctx context.Context, client WebSocketClient, endpoint string, callback func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, endpoint, callback)
}

func subscribeBytes(ctx context.Context, client WebSocketClient, endpoint string, callback func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	fullURL := client.GetHttpClient().BaseURL + endpoint
	dialer := client.GetDialer()
	conn, _, err := dialer.DialContext(ctx, fullURL, nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(655350)
	// doneC is buffered so a caller that sends on done instead of closing it
	// does not block once the connection has already ended.
	doneC := make(chan struct{}, 1)
	stopC := make(chan struct{})

	// Each read, and every ping or pong received while it waits, pushes the
	// read deadline back by timeout. The server answers the pings keepAlive
	// sends, so only a peer that has gone silent lets the deadline pass,
	// which fails ReadMessage and ends the subscription; time spent in
	// callback does not count. The handlers are set before reading starts,
	// since they may not change while ReadMessage runs. A failed pong write
	// is left to surface as a read error, so a close frame already received
	// is still reported.
	interval, timeout := keepAliveInterval, keepAliveTimeout
	extendDeadline := func() { conn.SetReadDeadline(time.Now().Add(timeout)) }
	conn.SetPongHandler(func(string) error {
		extendDeadline()
		return nil
	})
	conn.SetPingHandler(func(data string) error {
		extendDeadline()
		_ = conn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(controlWriteWait))
		return nil
	})
	go keepAlive(conn, interval, stopC)

	// closing marks a caller-initiated close, whose read error is expected
	// and not reported.
	var closing atomic.Bool
	go func() {
		select {
		case <-doneC:
			closing.Store(true)
		case <-stopC:
		}
		conn.Close()
	}()
	go func() {
		defer close(stopC)
		for {
			extendDeadline()
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !closing.Load() {
					callback(nil, err)
				}
				return
			}
			client.GetLogger().Debugf("received message: %s", common.BytesToString(message))
			callback(message, nil)
		}
	}()
	return doneC, stopC, nil
}

// keepAlive pings the server every interval until stop is closed.
func keepAlive(conn *websocket.Conn, interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
		if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(controlWriteWait)); err != nil {
			return
		}
	}
}
