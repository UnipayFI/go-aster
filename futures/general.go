package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type PingService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewPingService() *PingService {
	return &PingService{c: c}
}

func (s *PingService) Ping(ctx context.Context) error {
	req := request.Get(ctx, s.c, "/fapi/v1/ping")
	_, err := request.Do[struct{}](req)
	return err
}

type GetServerTimeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/time")
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime time.Time `json:"serverTime,format:unixmilli"`
}

type GetExchangeInfoService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetExchangeInfoService() *GetExchangeInfoService {
	return &GetExchangeInfoService{c: c}
}

func (s *GetExchangeInfoService) Do(ctx context.Context) (*ExchangeInfoResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/exchangeInfo")
	return request.Do[ExchangeInfoResponse](req)
}

type ExchangeInfoResponse struct {
	Timezone        string          `json:"timezone"`
	ServerTime      time.Time       `json:"serverTime,format:unixmilli"`
	RateLimits      []RateLimit     `json:"rateLimits"`
	ExchangeFilters []any           `json:"exchangeFilters"`
	Assets          []AssetInfo     `json:"assets"`
	Symbols         []FuturesSymbol `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

type AssetInfo struct {
	Asset             string          `json:"asset"`
	MarginAvailable   bool            `json:"marginAvailable"`
	AutoAssetExchange decimal.Decimal `json:"autoAssetExchange"`
}

type FuturesSymbol struct {
	Symbol                string           `json:"symbol"`
	Pair                  string           `json:"pair"`
	ContractType          string           `json:"contractType"`
	DeliveryDate          int64            `json:"deliveryDate"`
	OnboardDate           int64            `json:"onboardDate"`
	Status                string           `json:"status"`
	MaintMarginPercent    decimal.Decimal  `json:"maintMarginPercent"`
	RequiredMarginPercent decimal.Decimal  `json:"requiredMarginPercent"`
	BaseAsset             string           `json:"baseAsset"`
	QuoteAsset            string           `json:"quoteAsset"`
	MarginAsset           string           `json:"marginAsset"`
	PricePrecision        int              `json:"pricePrecision"`
	QuantityPrecision     int              `json:"quantityPrecision"`
	BaseAssetPrecision    int              `json:"baseAssetPrecision"`
	QuotePrecision        int              `json:"quotePrecision"`
	UnderlyingType        string           `json:"underlyingType"`
	UnderlyingSubType     []string         `json:"underlyingSubType"`
	SettlePlan            int              `json:"settlePlan"`
	TriggerProtect        decimal.Decimal  `json:"triggerProtect"`
	Filters               []map[string]any `json:"filters"`
	OrderTypes            []OrderType      `json:"OrderType"`
	TimeInForce           []TimeInForce    `json:"timeInForce"`
	LiquidationFee        decimal.Decimal  `json:"liquidationFee"`
	MarketTakeBound       decimal.Decimal  `json:"marketTakeBound"`
}
