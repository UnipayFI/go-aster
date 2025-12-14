package client

import (
	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-resty/resty/v2"
)

var _ request.Client = (*Client)(nil)

type Client struct {
	apiKey       string
	secretKey    string
	recvWindow   int64
	client       *resty.Client
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64
}

type SignFn = func(apiKey, secretKey string, payload string) (string, error)

func NewClient(options ...Options) *Client {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}
	return &Client{
		apiKey:       opt.apiKey,
		secretKey:    opt.secretKey,
		recvWindow:   opt.recvWindow,
		client:       opt.client,
		logger:       opt.logger,
		signFn:       opt.signFn,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *Client) GetHttpClient() *resty.Client {
	return c.client
}

func (c *Client) GetApiKey() string {
	return c.apiKey
}

func (c *Client) GetSecretKey() string {
	return c.secretKey
}

func (c *Client) GetRecvWindow() int64 {
	return c.recvWindow
}

func (c *Client) GetLogger() log.Logger {
	return c.logger
}

func (c *Client) GetSignFn() SignFn {
	return c.signFn
}

func (c *Client) GetTimeOffsetMs() int64 {
	return c.timeOffsetMs
}

// SetTimeOffset sets the time offset manually
func (c *Client) SetTimeOffset(offsetMs int64) {
	c.timeOffsetMs = offsetMs
}
