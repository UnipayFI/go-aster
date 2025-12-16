package client

import (
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-resty/resty/v2"
)

type WebSocketClient struct {
	apiKey       string
	secretKey    string
	client       *resty.Client
	recvWindow   int64
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64
}

func NewWebSocketClient(options ...Options) *WebSocketClient {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}
	return &WebSocketClient{
		apiKey:       opt.apiKey,
		secretKey:    opt.secretKey,
		recvWindow:   opt.recvWindow,
		client:       opt.client,
		logger:       opt.logger,
		signFn:       opt.signFn,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *WebSocketClient) GetApiKey() string {
	return c.apiKey
}

func (c *WebSocketClient) GetSecretKey() string {
	return c.secretKey
}

func (c *WebSocketClient) GetHttpClient() *resty.Client {
	return c.client
}

func (c *WebSocketClient) GetRecvWindow() int64 {
	return c.recvWindow
}

func (c *WebSocketClient) GetLogger() log.Logger {
	return c.logger
}

func (c *WebSocketClient) GetSignFn() SignFn {
	return c.signFn
}

func (c *WebSocketClient) GetTimeOffsetMs() int64 {
	return c.timeOffsetMs
}

// SetTimeOffset sets the time offset manually
func (c *WebSocketClient) SetTimeOffset(offsetMs int64) {
	c.timeOffsetMs = offsetMs
}

func (c *WebSocketClient) GetBaseURL() string {
	return c.client.BaseURL
}
