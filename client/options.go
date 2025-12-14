package client

import (
	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-json-experiment/json"
	"github.com/go-resty/resty/v2"
)

const (
	DEFAULT_RECV_WINDOW = 5000
)

type Option struct {
	apiKey       string
	secretKey    string
	recvWindow   int64
	logger       log.Logger
	signFn       SignFn
	client       *resty.Client
	timeOffsetMs int64
}

type Options func(*Option)

func defaultOption() *Option {
	return &Option{
		recvWindow: DEFAULT_RECV_WINDOW,
		logger:     log.GetDefaultLogger(),
		client:     defaultHttpClient(),
	}
}

func defaultHttpClient() *resty.Client {
	return resty.New().
		SetBaseURL(common.DEFAULT_SPOT_BASE_URL).
		SetJSONMarshaler(func(v any) ([]byte, error) {
			return json.Marshal(v)
		}).
		SetJSONUnmarshaler(func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		})
}

func WithLogger(logger log.Logger) Options {
	return func(opt *Option) {
		opt.logger = logger
	}
}

func WithAuth(apiKey, secretKey string) Options {
	return func(opt *Option) {
		opt.apiKey = apiKey
		opt.secretKey = secretKey
	}
}

// WithSignFn sets the sign function for the client.
func WithSignRequestFn(signFn SignFn) Options {
	return func(opt *Option) {
		opt.signFn = signFn
	}
}

func WithTimeOffset(timeOffsetMs int64) Options {
	return func(opt *Option) {
		opt.timeOffsetMs = timeOffsetMs
	}
}

func WithBaseURL(baseURL string) Options {
	return func(opt *Option) {
		opt.client.SetBaseURL(baseURL)
	}
}
