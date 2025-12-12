package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type DepthService struct {
	client *SpotClient

	params map[string]string
}

func NewDepthService(client *SpotClient, symbol string) *DepthService {
	return &DepthService{client: client, params: map[string]string{"symbol": symbol}}
}

func (s *DepthService) SetLimit(limit int) *DepthService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *DepthService) SetSymbolStatus(symbolStatus string) *DepthService {
	s.params["symbolStatus"] = symbolStatus
	return s
}

func (s *DepthService) Do(ctx context.Context) (*DepthResponse, error) {
	req := request.Get(ctx, s.client, "/api/v1/depth", s.params)
	return request.Do[DepthResponse](req)
}

type DepthResponse struct {
	LastUpdateId int64                `json:"lastUpdateId"`
	Bids         [][2]decimal.Decimal `json:"bids"`
	Asks         [][2]decimal.Decimal `json:"asks"`
}

type GetTradesService struct {
	client *SpotClient

	params map[string]string
}

func NewGetTradesService(client *SpotClient, symbol string) *GetTradesService {
	return &GetTradesService{client: client, params: map[string]string{"symbol": symbol}}
}

func (s *GetTradesService) SetLimit(limit int) *GetTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetTradesService) Do(ctx context.Context) ([]TradeResponse, error) {
	req := request.Get(ctx, s.client, "/api/v1/trades", s.params)
	trades, err := request.Do[[]TradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type TradeResponse struct {
	Id           int64           `json:"id"`
	Price        decimal.Decimal `json:"price"`
	Quantity     decimal.Decimal `json:"qty"`
	BaseQuantity decimal.Decimal `json:"baseQty"`
	Time         time.Time       `json:"time,format:unixmilli"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
}
