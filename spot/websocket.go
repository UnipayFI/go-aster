package spot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/common"
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

func (s *SubscribeAggTradeService) Do(ctx context.Context, handler func(message *WsAggTradeResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@aggTrade"
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

type SubscribeCombinedDepthService struct {
	c            *SpotWebSocketClient
	symbolLevels map[string]string
}

// symbolLevels:
// "BTCUSDT": "5@100ms",
// "ETHUSDT": "5@100ms",
// "SOLUSDT": "5@100ms",
func (c *SpotWebSocketClient) NewSubscribeCombinedDepthService(symbolLevels map[string]string) *SubscribeCombinedDepthService {
	return &SubscribeCombinedDepthService{c: c, symbolLevels: symbolLevels}
}

func (s *SubscribeCombinedDepthService) Do(ctx context.Context, handler func(message *WsCombinedDepthResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	symbols := make([]string, 0, len(s.symbolLevels))
	for symbol, level := range s.symbolLevels {
		symbols = append(symbols, fmt.Sprintf("%s@depth%s", strings.ToLower(symbol), level))
	}
	url := "/stream?streams=" + strings.Join(symbols, "/")
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsCombinedDepthResponse struct {
	Stream string          `json:"stream"`
	Data   WsDepthResponse `json:"data"`
}

type WsDepthResponse struct {
	WsBaseEvent
	TransactionTime  time.Time   `json:"T,format:unixmilli"`
	Symbol           string      `json:"s"`
	FirstUpdateID    int64       `json:"U"`
	LastUpdateID     int64       `json:"u"`
	PrevLastUpdateID int64       `json:"pu"`
	Bids             []PriceSize `json:"b"`
	Asks             []PriceSize `json:"a"`
}
