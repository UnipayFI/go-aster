package client

import (
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client

	apiKey       string
	secretKey    string
	recvWindow   int64
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64

	UsedWeight UsedWeight
	OrderCount OrderCount
}

type UsedWeight struct {
	Used   int64
	Used1M int64 // used in last 1 minute
}

type OrderCount struct {
	Count10s int64
	Count1d  int64
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

func (c *Client) GetUsedWeight() UsedWeight {
	return c.UsedWeight
}

func (c *Client) GetOrderCount() OrderCount {
	return c.OrderCount
}

func (c *Client) SetUsedWeight(used, used1m *int64) {
	if used != nil {
		c.UsedWeight.Used = *used
	}
	if used1m != nil {
		c.UsedWeight.Used1M = *used1m
	}
}

func (c *Client) SetOrderCount(count10s, count1d *int64) {
	if count10s != nil {
		c.OrderCount.Count10s = *count10s
	}
	if count1d != nil {
		c.OrderCount.Count1d = *count1d
	}
}
