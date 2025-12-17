package spot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
)

type SubscribeAggTradeService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeAggTradeService(symbol string) *SubscribeAggTradeService {
	return &SubscribeAggTradeService{c: c, symbol: symbol}
}

func (s *SubscribeAggTradeService) Do(ctx context.Context, handler func(message *WsAggTradeResponse, err error)) (doneC <-chan struct{}, stopC chan struct{}, err error) {
	url := fmt.Sprintf("%s@aggTrade", strings.ToLower(s.symbol))
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsAggTradeResponse struct {
	EventType    string          `json:"e"`
	EventTime    time.Time       `json:"E,format:unixmilli"`
	Symbol       string          `json:"s"`
	AggTradeId   int64           `json:"a"`
	Price        decimal.Decimal `json:"p"`
	Quantity     decimal.Decimal `json:"q"`
	FirstTradeId int64           `json:"f"`
	LastTradeId  int64           `json:"l"`
	TradeTime    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}
