package client

import (
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	client *resty.Client
	logger log.Logger
	dialer *websocket.Dialer
}

func NewWebSocketClient(options ...WebSocketOptions) *WebSocketClient {
	opt := defaultWebSocketOption()
	for _, option := range options {
		option(opt)
	}
	return &WebSocketClient{
		client: opt.client,
		logger: opt.logger,
		dialer: opt.dialer,
	}
}

func (c *WebSocketClient) GetHttpClient() *resty.Client {
	return c.client
}

func (c *WebSocketClient) GetLogger() log.Logger {
	return c.logger
}

func (c *WebSocketClient) GetBaseURL() string {
	return c.client.BaseURL
}

func (c *WebSocketClient) GetDialer() *websocket.Dialer {
	return c.dialer
}
