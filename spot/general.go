package spot

import (
	"context"
	"encoding/json"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
)

type PingService struct {
	c *SpotClient
}

func (c *SpotClient) NewPingService() *PingService {
	return &PingService{c: c}
}

func (s *PingService) Ping(ctx context.Context) error {
	req := request.Get(ctx, s.c, "/api/v1/ping")
	_, err := request.Do[json.RawMessage](req)
	if err != nil {
		return err
	}
	return nil
}

type TimeService struct {
	c *SpotClient
}

func (c *SpotClient) NewTimeService() *TimeService {
	return &TimeService{c: c}
}

func (s *TimeService) GetTime(ctx context.Context) (*TimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/time")
	return request.Do[TimeResponse](req)
}

type TimeResponse struct {
	ServerTime time.Time `json:"serverTime,format:unixmilli"`
}

type ExchangeInfoService struct {
	c *SpotClient
}

func (c *SpotClient) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

func (s *ExchangeInfoService) GetExchangeInfo(ctx context.Context) (*ExchangeInfoResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/exchangeInfo")
	return request.Do[ExchangeInfoResponse](req)
}

type ExchangeInfoResponse struct {
	Timezone        string      `json:"timezone"`
	ServerTime      time.Time   `json:"serverTime,format:unixmilli"`
	RateLimits      []RateLimit `json:"rateLimits"`
	ExchangeFilters []any       `json:"exchangeFilters"`
	Assets          []AssetInfo `json:"assets"`
	Symbols         []Symbol    `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

type AssetInfo struct {
	Asset string `json:"asset"`
}

type Symbol struct {
	Symbol             string           `json:"symbol"`
	Status             string           `json:"status"`
	BaseAsset          string           `json:"baseAsset"`
	QuoteAsset         string           `json:"quoteAsset"`
	PricePrecision     int              `json:"pricePrecision"`
	QuantityPrecision  int              `json:"quantityPrecision"`
	BaseAssetPrecision int              `json:"baseAssetPrecision"`
	QuotePrecision     int              `json:"quotePrecision"`
	ListingTime        time.Time        `json:"listingTime,format:unixmilli"`
	BaseAssetAddress   string           `json:"baseAssetAddress"`
	Filters            []map[string]any `json:"filters"`
	OrderTypes         []OrderType      `json:"orderTypes"`
	TimeInForce        []TimeInForce    `json:"timeInForce"`
	OcoAllowed         bool             `json:"ocoAllowed"`
}

type GetServerTimeService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/time", map[string]string{})
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}
