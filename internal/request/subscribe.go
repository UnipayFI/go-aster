package request

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/gorilla/websocket"
)

func Subscribe[T any](ctx context.Context, client Client, endpoint string, callback func(message *T, err error)) (doneC <-chan struct{}, stopC chan struct{}, err error) {
	fullURL, _ := url.JoinPath(client.GetHttpClient().BaseURL, endpoint)
	dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, err := dialer.DialContext(ctx, fullURL, nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})

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
			var msg T
			err = client.GetHttpClient().JSONUnmarshal(message, &msg)
			callback(&msg, err)
		}
	}()
	return doneC, stopC, nil
}

func keepAlive(conn *websocket.Conn, timeout, interval time.Duration) {
	latest := time.Now()

	// Set pong handler to keep the connection alive
	conn.SetPongHandler(func(raw string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte(raw), time.Now().Add(timeout))
		if err == nil {
			latest = time.Now()
		}
		return err
	})
	// Set ping handler to keep the connection alive
	conn.SetPingHandler(func(raw string) error {
		err := conn.WriteControl(websocket.PingMessage, []byte(raw), time.Now().Add(timeout))
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
