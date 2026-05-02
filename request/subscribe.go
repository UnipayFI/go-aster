package request

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

type WebSocketClient interface {
	GetHttpClient() *resty.Client
	GetLogger() log.Logger
	GetDialer() *websocket.Dialer
}

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
	doneC := make(chan struct{})
	stopC := make(chan struct{})

	go keepAlive(conn, common.DEFAULT_KEEP_ALIVE_TIMEOUT, common.DEFAULT_KEEP_ALIVE_INTERVAL)

	silent := false
	go func() {
		select {
		case <-stopC:
			silent = true
		case <-doneC:
		}
		conn.Close()
	}()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !silent {
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

func keepAlive(conn *websocket.Conn, timeout, interval time.Duration) {
	latest := time.Now()

	// Update activity timestamp on incoming pongs (no reply needed for pongs).
	conn.SetPongHandler(func(string) error {
		latest = time.Now()
		return nil
	})
	// Reply to incoming pings with a pong (the gorilla default does this too,
	// but we override to also update the activity timestamp).
	conn.SetPingHandler(func(raw string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte(raw), time.Now().Add(timeout))
		if err == nil {
			latest = time.Now()
		}
		return err
	})

	// Set ticker to send ping messages at the interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		if time.Since(latest) > timeout {
			conn.Close()
			return
		}
		err := conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			return
		}
		latest = time.Now()
	}
}
