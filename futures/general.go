package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
)

// PingService -- GET /fapi/v3/ping
type PingService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewPingService() *PingService {
	return &PingService{c: c}
}

func (s *PingService) Do(ctx context.Context) error {
	req := request.Get(ctx, s.c, "/fapi/v3/ping")
	_, err := request.Do[struct{}](req)
	return err
}

// GetServerTimeService -- GET /fapi/v3/time
type GetServerTimeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/time")
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime time.Time `json:"serverTime,format:unixmilli"`
}

// NoopService -- POST /fapi/v3/noop
//
// V3 fast nonce-bump helper for market makers; requires TRADE auth.
type NoopService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewNoopService() *NoopService {
	return &NoopService{c: c}
}

func (s *NoopService) Do(ctx context.Context) (*NoopResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/noop").WithSignature()
	return request.Do[NoopResponse](req)
}

type NoopResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
