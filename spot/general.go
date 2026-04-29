package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
)

// PingService -- GET /api/v3/ping
//
// Test connectivity to the REST API. Returns no body on success.
type PingService struct {
	c *SpotClient
}

func (c *SpotClient) NewPingService() *PingService {
	return &PingService{c: c}
}

func (s *PingService) Do(ctx context.Context) error {
	req := request.Get(ctx, s.c, "/api/v3/ping")
	_, err := request.Do[struct{}](req)
	return err
}

// GetServerTimeService -- GET /api/v3/time
type GetServerTimeService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/time")
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime time.Time `json:"serverTime,format:unixmilli"`
}

// NoopService -- POST /api/v3/noop
//
// "Noop" is a V3-specific fast cancel-helper endpoint: it accepts a signed
// request, increments the user's nonce floor, and returns success. Used by
// market makers to invalidate stale nonces ahead of a real order. Requires
// SPOT_TRADE auth.
type NoopService struct {
	c *SpotClient
}

func (c *SpotClient) NewNoopService() *NoopService {
	return &NoopService{c: c}
}

func (s *NoopService) Do(ctx context.Context) (*NoopResponse, error) {
	req := request.Post(s.c, ctx, "/api/v3/noop").WithSignature()
	return request.Do[NoopResponse](req)
}

type NoopResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
