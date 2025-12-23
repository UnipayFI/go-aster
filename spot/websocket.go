package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
)

type DiffDepthSpeed string

const (
	// DiffDepthSpeed100ms represents the `@100ms` suffix.
	DiffDepthSpeed100ms DiffDepthSpeed = "100ms"
)

type SubscribeAggTradeService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeAggTradeService(symbol string) *SubscribeAggTradeService {
	return &SubscribeAggTradeService{c: c, symbol: symbol}
}

func (s *SubscribeAggTradeService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "aggTrade")
}

func (s *SubscribeAggTradeService) Do(ctx context.Context, handler func(message *WsAggTradeResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
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

type SubscribeTradeService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeTradeService(symbol string) *SubscribeTradeService {
	return &SubscribeTradeService{c: c, symbol: symbol}
}

func (s *SubscribeTradeService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "trade")
}

func (s *SubscribeTradeService) Do(ctx context.Context, handler func(message *WsTradeResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type WsTradeResponse struct {
	EventType    string          `json:"e"`
	EventTime    time.Time       `json:"E,format:unixmilli"`
	Symbol       string          `json:"s"`
	TradeId      int64           `json:"t"`
	Price        decimal.Decimal `json:"p"`
	Quantity     decimal.Decimal `json:"q"`
	TradeTime    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}

type SubscribeKlineService struct {
	c        *SpotWebSocketClient
	symbol   string
	interval KlineInterval
}

func (c *SpotWebSocketClient) NewSubscribeKlineService(symbol string, interval KlineInterval) *SubscribeKlineService {
	return &SubscribeKlineService{c: c, symbol: symbol, interval: interval}
}

func (s *SubscribeKlineService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "kline_"+string(s.interval))
}

func (s *SubscribeKlineService) Do(ctx context.Context, handler func(message *WsKlineResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type WsKlineResponse struct {
	EventType string    `json:"e"`
	EventTime time.Time `json:"E,format:unixmilli"`
	Symbol    string    `json:"s"`
	Kline     WsKline   `json:"k"`
}

type WsKline struct {
	StartTime           time.Time       `json:"t,format:unixmilli"`
	CloseTime           time.Time       `json:"T,format:unixmilli"`
	Symbol              string          `json:"s"`
	Interval            KlineInterval   `json:"i"`
	FirstTradeId        int64           `json:"f"`
	LastTradeId         int64           `json:"L"`
	Open                decimal.Decimal `json:"o"`
	Close               decimal.Decimal `json:"c"`
	High                decimal.Decimal `json:"h"`
	Low                 decimal.Decimal `json:"l"`
	Volume              decimal.Decimal `json:"v"`
	TradeCount          int64           `json:"n"`
	IsClosed            bool            `json:"x"`
	QuoteVolume         decimal.Decimal `json:"q"`
	TakerBuyBaseVolume  decimal.Decimal `json:"V"`
	TakerBuyQuoteVolume decimal.Decimal `json:"Q"`
}

type SubscribeMiniTickerService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeMiniTickerService(symbol string) *SubscribeMiniTickerService {
	return &SubscribeMiniTickerService{c: c, symbol: symbol}
}

func (s *SubscribeMiniTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "miniTicker")
}

func (s *SubscribeMiniTickerService) Do(ctx context.Context, handler func(message *WsMiniTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllMiniTickersService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllMiniTickersService() *SubscribeAllMiniTickersService {
	return &SubscribeAllMiniTickersService{c: c}
}

func (s *SubscribeAllMiniTickersService) endpoint() string {
	return common.WsAllStream("miniTicker@arr")
}

func (s *SubscribeAllMiniTickersService) Do(ctx context.Context, handler func(message []*WsMiniTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := s.endpoint()
	callback := func(message *[]*WsMiniTickerResponse, err error) {
		if message == nil {
			handler(nil, err)
			return
		}
		handler(*message, err)
	}
	return request.Subscribe(ctx, s.c, url, callback)
}

type WsMiniTickerResponse struct {
	EventType   string          `json:"e"`
	EventTime   time.Time       `json:"E,format:unixmilli"`
	Symbol      string          `json:"s"`
	ClosePrice  decimal.Decimal `json:"c"`
	OpenPrice   decimal.Decimal `json:"o"`
	HighPrice   decimal.Decimal `json:"h"`
	LowPrice    decimal.Decimal `json:"l"`
	Volume      decimal.Decimal `json:"v"`
	QuoteVolume decimal.Decimal `json:"q"`
}

type SubscribeTickerService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeTickerService(symbol string) *SubscribeTickerService {
	return &SubscribeTickerService{c: c, symbol: symbol}
}

func (s *SubscribeTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "ticker")
}

func (s *SubscribeTickerService) Do(ctx context.Context, handler func(message *WsTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllTickersService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllTickersService() *SubscribeAllTickersService {
	return &SubscribeAllTickersService{c: c}
}

func (s *SubscribeAllTickersService) endpoint() string {
	return common.WsAllStream("ticker@arr")
}

func (s *SubscribeAllTickersService) Do(ctx context.Context, handler func(message []*WsTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := s.endpoint()
	callback := func(message *[]*WsTickerResponse, err error) {
		if message == nil {
			handler(nil, err)
			return
		}
		handler(*message, err)
	}
	return request.Subscribe(ctx, s.c, url, callback)
}

type WsTickerResponse struct {
	EventType          string          `json:"e"`
	EventTime          time.Time       `json:"E,format:unixmilli"`
	Symbol             string          `json:"s"`
	PriceChange        decimal.Decimal `json:"p"`
	PriceChangePercent decimal.Decimal `json:"P"`
	WeightedAvgPrice   decimal.Decimal `json:"w"`
	LastPrice          decimal.Decimal `json:"c"`
	LastQuantity       decimal.Decimal `json:"Q"`
	OpenPrice          decimal.Decimal `json:"o"`
	HighPrice          decimal.Decimal `json:"h"`
	LowPrice           decimal.Decimal `json:"l"`
	Volume             decimal.Decimal `json:"v"`
	QuoteVolume        decimal.Decimal `json:"q"`
	OpenTime           time.Time       `json:"O,format:unixmilli"`
	CloseTime          time.Time       `json:"C,format:unixmilli"`
	FirstTradeId       int64           `json:"F"`
	LastTradeId        int64           `json:"L"`
	TradeCount         int64           `json:"n"`
}

type SubscribeBookTickerService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeBookTickerService(symbol string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, symbol: symbol}
}

func (s *SubscribeBookTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "bookTicker")
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllBookTickersService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllBookTickersService() *SubscribeAllBookTickersService {
	return &SubscribeAllBookTickersService{c: c}
}

func (s *SubscribeAllBookTickersService) endpoint() string {
	return common.WsAllStream("bookTicker")
}

func (s *SubscribeAllBookTickersService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type WsBookTickerResponse struct {
	UpdateId     int64           `json:"u"`
	Symbol       string          `json:"s"`
	BestBidPrice decimal.Decimal `json:"b"`
	BestBidQty   decimal.Decimal `json:"B"`
	BestAskPrice decimal.Decimal `json:"a"`
	BestAskQty   decimal.Decimal `json:"A"`
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

func (s *SubscribeCombinedDepthService) endpoint() string {
	streams := make([]string, 0, len(s.symbolLevels))
	for symbol, level := range s.symbolLevels {
		streams = append(streams, common.WsSymbolStreamSegment(symbol, "depth"+level))
	}
	return common.WsCombinedStreamsEndpoint(streams...)
}

func (s *SubscribeCombinedDepthService) Do(ctx context.Context, handler func(message *WsCombinedDepthResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
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

type SubscribeDiffDepthService struct {
	c      *SpotWebSocketClient
	symbol string
	speed  DiffDepthSpeed
}

func (c *SpotWebSocketClient) NewSubscribeDiffDepthService(symbol string) *SubscribeDiffDepthService {
	return &SubscribeDiffDepthService{c: c, symbol: symbol}
}

// Default update frequency is 1000ms per docs when Speed is not set.
func (s *SubscribeDiffDepthService) Speed(speed DiffDepthSpeed) *SubscribeDiffDepthService {
	s.speed = speed
	return s
}

func (s *SubscribeDiffDepthService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "depth", string(s.speed))
}

func (s *SubscribeDiffDepthService) Do(ctx context.Context, handler func(message *WsDepthResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}
