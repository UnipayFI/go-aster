package spot

import (
	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-resty/resty/v2"
)

var _ request.Client = (*SpotClient)(nil)

type SpotClient struct {
	apiKey     string
	secretKey  string
	recvWindow int64
	client     *resty.Client
	logger     log.Logger
	signFn     SignFn
}

type SignFn = func(apiKey, secretKey string, payload string) (string, error)

func NewSpotClient(options ...Options) *SpotClient {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}
	return &SpotClient{
		apiKey:     opt.apiKey,
		secretKey:  opt.secretKey,
		recvWindow: opt.recvWindow,
		client:     opt.client,
		logger:     opt.logger,
		signFn:     opt.signFn,
	}
}

func (c *SpotClient) GetHttpClient() *resty.Client {
	return c.client
}

func (c *SpotClient) GetApiKey() string {
	return c.apiKey
}

func (c *SpotClient) GetSecretKey() string {
	return c.secretKey
}

func (c *SpotClient) GetRecvWindow() int64 {
	return c.recvWindow
}

func (c *SpotClient) GetLogger() log.Logger {
	return c.logger
}

func (c *SpotClient) GetSignFn() SignFn {
	return c.signFn
}
