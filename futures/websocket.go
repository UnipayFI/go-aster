package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
)

type MarkPriceSpeed string

const (
	// MarkPriceSpeed1s represents the `@1s` suffix.
	MarkPriceSpeed1s MarkPriceSpeed = "1s"
)

type DiffDepthSpeed string

const (
	DiffDepthSpeed100ms DiffDepthSpeed = "100ms"
	DiffDepthSpeed500ms DiffDepthSpeed = "500ms"
)

type SubscribeAggTradeService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeAggTradeService(symbol string) *SubscribeAggTradeService {
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

type SubscribeMarkPriceService struct {
	c      *FuturesWebSocketClient
	symbol string
	speed  MarkPriceSpeed
}

func (c *FuturesWebSocketClient) NewSubscribeMarkPriceService(symbol string) *SubscribeMarkPriceService {
	return &SubscribeMarkPriceService{c: c, symbol: symbol}
}

// Default update frequency is 3000ms per docs when Speed is not set.
func (s *SubscribeMarkPriceService) Speed(speed MarkPriceSpeed) *SubscribeMarkPriceService {
	s.speed = speed
	return s
}

func (s *SubscribeMarkPriceService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "markPrice", string(s.speed))
}

func (s *SubscribeMarkPriceService) Do(ctx context.Context, handler func(message *WsMarkPriceResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllMarkPricesService struct {
	c     *FuturesWebSocketClient
	speed MarkPriceSpeed
}

func (c *FuturesWebSocketClient) NewSubscribeAllMarkPricesService() *SubscribeAllMarkPricesService {
	return &SubscribeAllMarkPricesService{c: c}
}

// Default update frequency is 3000ms per docs when Speed is not set.
func (s *SubscribeAllMarkPricesService) Speed(speed MarkPriceSpeed) *SubscribeAllMarkPricesService {
	s.speed = speed
	return s
}

func (s *SubscribeAllMarkPricesService) endpoint() string {
	return common.WsAllStream("markPrice@arr", string(s.speed))
}

func (s *SubscribeAllMarkPricesService) Do(ctx context.Context, handler func(message []*WsMarkPriceResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := s.endpoint()
	callback := func(message *[]*WsMarkPriceResponse, err error) {
		if message == nil {
			handler(nil, err)
			return
		}
		handler(*message, err)
	}
	return request.Subscribe(ctx, s.c, url, callback)
}

type WsMarkPriceResponse struct {
	EventType       string          `json:"e"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	Symbol          string          `json:"s"`
	MarkPrice       decimal.Decimal `json:"p"`
	IndexPrice      decimal.Decimal `json:"i"`
	EstSettlePrice  decimal.Decimal `json:"P"`
	FundingRate     decimal.Decimal `json:"r"`
	NextFundingTime time.Time       `json:"T,format:unixmilli"`
}

type SubscribeKlineService struct {
	c        *FuturesWebSocketClient
	symbol   string
	interval KlineInterval
}

func (c *FuturesWebSocketClient) NewSubscribeKlineService(symbol string, interval KlineInterval) *SubscribeKlineService {
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
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeMiniTickerService(symbol string) *SubscribeMiniTickerService {
	return &SubscribeMiniTickerService{c: c, symbol: symbol}
}

func (s *SubscribeMiniTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "miniTicker")
}

func (s *SubscribeMiniTickerService) Do(ctx context.Context, handler func(message *WsMiniTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllMiniTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllMiniTickersService() *SubscribeAllMiniTickersService {
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
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeTickerService(symbol string) *SubscribeTickerService {
	return &SubscribeTickerService{c: c, symbol: symbol}
}

func (s *SubscribeTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "ticker")
}

func (s *SubscribeTickerService) Do(ctx context.Context, handler func(message *WsTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllTickersService() *SubscribeAllTickersService {
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
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeBookTickerService(symbol string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, symbol: symbol}
}

func (s *SubscribeBookTickerService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "bookTicker")
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllBookTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllBookTickersService() *SubscribeAllBookTickersService {
	return &SubscribeAllBookTickersService{c: c}
}

func (s *SubscribeAllBookTickersService) endpoint() string {
	return common.WsAllStream("bookTicker")
}

func (s *SubscribeAllBookTickersService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type WsBookTickerResponse struct {
	EventType       string          `json:"e"`
	UpdateId        int64           `json:"u"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	TransactionTime time.Time       `json:"T,format:unixmilli"`
	Symbol          string          `json:"s"`
	BestBidPrice    decimal.Decimal `json:"b"`
	BestBidQty      decimal.Decimal `json:"B"`
	BestAskPrice    decimal.Decimal `json:"a"`
	BestAskQty      decimal.Decimal `json:"A"`
}

type SubscribeCombinedDepthService struct {
	c            *FuturesWebSocketClient
	symbolLevels map[string]string
}

// symbolLevels:
// "BTCUSDT": "5@100ms",
// "ETHUSDT": "10@100ms",
func (c *FuturesWebSocketClient) NewSubscribeCombinedDepthService(symbolLevels map[string]string) *SubscribeCombinedDepthService {
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
	c      *FuturesWebSocketClient
	symbol string
	speed  DiffDepthSpeed
}

func (c *FuturesWebSocketClient) NewSubscribeDiffDepthService(symbol string) *SubscribeDiffDepthService {
	return &SubscribeDiffDepthService{c: c, symbol: symbol}
}

// Default update frequency is 250ms per docs when Speed is not set.
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

type SubscribeForceOrderService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeForceOrderService(symbol string) *SubscribeForceOrderService {
	return &SubscribeForceOrderService{c: c, symbol: symbol}
}

func (s *SubscribeForceOrderService) endpoint() string {
	return common.WsSymbolStream(s.symbol, "forceOrder")
}

func (s *SubscribeForceOrderService) Do(ctx context.Context, handler func(message *WsForceOrderResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.Subscribe(ctx, s.c, s.endpoint(), handler)
}

type SubscribeAllForceOrdersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllForceOrdersService() *SubscribeAllForceOrdersService {
	return &SubscribeAllForceOrdersService{c: c}
}

func (s *SubscribeAllForceOrdersService) endpoint() string {
	return common.WsAllStream("forceOrder@arr")
}

func (s *SubscribeAllForceOrdersService) Do(ctx context.Context, handler func(message []*WsForceOrderResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := s.endpoint()
	callback := func(message *[]*WsForceOrderResponse, err error) {
		if message == nil {
			handler(nil, err)
			return
		}
		handler(*message, err)
	}
	return request.Subscribe(ctx, s.c, url, callback)
}

type WsForceOrderResponse struct {
	EventType string        `json:"e"`
	EventTime time.Time     `json:"E,format:unixmilli"`
	Order     WsLiquidation `json:"o"`
}

type WsLiquidation struct {
	Symbol           string          `json:"s"`
	Side             OrderSide       `json:"S"`
	OrderType        OrderType       `json:"o"`
	TimeInForce      TimeInForce     `json:"f"`
	OriginalQuantity decimal.Decimal `json:"q"`
	Price            decimal.Decimal `json:"p"`
	AveragePrice     decimal.Decimal `json:"ap"`
	Status           OrderStatus     `json:"X"`
	LastFilledQty    decimal.Decimal `json:"l"`
	AccumulatedQty   decimal.Decimal `json:"z"`
	TradeTime        time.Time       `json:"T,format:unixmilli"`
}
