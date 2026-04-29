package client

import (
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-resty/resty/v2"
)

type WebSocketClient struct {
	client *resty.Client
	logger log.Logger
}

func NewWebSocketClient(options ...WebSocketOptions) *WebSocketClient {
	opt := defaultWebSocketOption()
	for _, option := range options {
		option(opt)
	}
	return &WebSocketClient{
		client: opt.client,
		logger: opt.logger,
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
