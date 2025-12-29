package futures

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
)

type subscribeAggTradeService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeAggTradeService(symbol string) *subscribeAggTradeService {
	return &subscribeAggTradeService{c: c, symbol: symbol}
}

func (s *subscribeAggTradeService) Do(ctx context.Context, handler func(message *WsAggTradeResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
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

type subscribeMarkPriceService struct {
	c      *FuturesWebSocketClient
	symbol string
	speed  string
}

// speed: 1s
func (c *FuturesWebSocketClient) NewSubscribeMarkPriceService(symbol string, speed ...string) *subscribeMarkPriceService {
	s := subscribeMarkPriceService{c: c, symbol: symbol}
	if len(speed) > 0 {
		s.speed = "@" + speed[0]
	}
	return &s
}

func (s *subscribeMarkPriceService) Do(ctx context.Context, handler func(message *WsMarkPriceResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@markPrice"
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsMarkPriceResponse struct {
	EventType       string          `json:"e"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	Symbol          string          `json:"s"`
	Price           decimal.Decimal `json:"p"`
	IndexPrice      decimal.Decimal `json:"i"`
	EstimatedPrice  decimal.Decimal `json:"P"`
	FundingRate     decimal.Decimal `json:"r"`
	NextFundingTime time.Time       `json:"T,format:unixmilli"`
}

type subscribeAllMarkPriceService struct {
	c     *FuturesWebSocketClient
	speed string
}

func (c *FuturesWebSocketClient) NewSubscribeAllMarkPriceService(speed ...string) *subscribeAllMarkPriceService {
	s := subscribeAllMarkPriceService{c: c}
	if len(speed) > 0 {
		s.speed = "@" + speed[0]
	}
	return &s
}

func (s *subscribeAllMarkPriceService) Do(ctx context.Context, handler func(message *[]WsMarkPriceResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + "!markPrice@arr" + s.speed
	return request.Subscribe(ctx, s.c, url, handler)
}

type subscribeKlineService struct {
	c        *FuturesWebSocketClient
	symbol   string
	interval KlineInterval
}

func (c *FuturesWebSocketClient) NewSubscribeKlineService(symbol string, interval KlineInterval) *subscribeKlineService {
	return &subscribeKlineService{c: c, symbol: symbol, interval: interval}
}

func (s *subscribeKlineService) Do(ctx context.Context, handler func(message *WsKlineResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@kline_" + string(s.interval)
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsKlineResponse struct {
	EventType string    `json:"e"`
	EventTime time.Time `json:"E,format:unixmilli"`
	Symbol    string    `json:"s"`
	Kline     WsKline   `json:"k"`
}

type WsKline struct {
	StartTime                time.Time       `json:"t,format:unixmilli"`
	EndTime                  time.Time       `json:"T,format:unixmilli"`
	Symbol                   string          `json:"s"`
	Interval                 KlineInterval   `json:"i"`
	FirstTradeId             int64           `json:"f"`
	LastTradeId              int64           `json:"L"`
	OpenPrice                decimal.Decimal `json:"o"`
	ClosePrice               decimal.Decimal `json:"c"`
	HighPrice                decimal.Decimal `json:"h"`
	LowPrice                 decimal.Decimal `json:"l"`
	Volume                   decimal.Decimal `json:"v"`
	TradeNum                 int64           `json:"n"`
	IsClosed                 bool            `json:"x"`
	QuoteAssetVolume         decimal.Decimal `json:"q"`
	TakerBuyBaseAssetVolume  decimal.Decimal `json:"V"`
	TakerBuyQuoteAssetVolume decimal.Decimal `json:"Q"`
	BuyerBaseAssetVolume     decimal.Decimal `json:"B"`
}

type subscribeMiniTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeMiniTickerService(symbol string) *subscribeMiniTickerService {
	return &subscribeMiniTickerService{c: c, symbol: symbol}
}

func (s *subscribeMiniTickerService) Do(ctx context.Context, handler func(message *WsMiniTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@miniTicker"
	return request.Subscribe(ctx, s.c, url, handler)
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

type subscribeAllMiniTickerService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllMiniTickerService() *subscribeAllMiniTickerService {
	return &subscribeAllMiniTickerService{c: c}
}

func (s *subscribeAllMiniTickerService) Do(ctx context.Context, handler func(message *[]WsMiniTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + "!miniTicker@arr"
	return request.Subscribe(ctx, s.c, url, handler)
}

type subscribeTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeTickerService(symbol string) *subscribeTickerService {
	return &subscribeTickerService{c: c, symbol: symbol}
}

func (s *subscribeTickerService) Do(ctx context.Context, handler func(message *WsTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@ticker"
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsTickerResponse struct {
	EventType           string          `json:"e"`
	EventTime           time.Time       `json:"E,format:unixmilli"`
	Symbol              string          `json:"s"`
	PriceChange         decimal.Decimal `json:"p"`
	PriceChangePercent  decimal.Decimal `json:"P"`
	WeightedAvgPrice    decimal.Decimal `json:"w"`
	LastPrice           decimal.Decimal `json:"c"`
	LastQuantity        decimal.Decimal `json:"Q"`
	OpenPrice           decimal.Decimal `json:"o"`
	HighPrice           decimal.Decimal `json:"h"`
	LowPrice            decimal.Decimal `json:"l"`
	Volume              decimal.Decimal `json:"v"`
	QuoteVolume         decimal.Decimal `json:"q"`
	StatisticsOpenTime  time.Time       `json:"O,format:unixmilli"`
	StatisticsCloseTime time.Time       `json:"C,format:unixmilli"`
	FirstTradeId        int64           `json:"F"`
	LastTradeId         int64           `json:"L"`
	TotalTrades         int64           `json:"n"`
}

type subscribeAllTickerService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllTickerService() *subscribeAllTickerService {
	return &subscribeAllTickerService{c: c}
}

func (s *subscribeAllTickerService) Do(ctx context.Context, handler func(message *[]WsTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + "!ticker@arr"
	return request.Subscribe(ctx, s.c, url, handler)
}

type subscribeBookTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeBookTickerService(symbol string) *subscribeBookTickerService {
	return &subscribeBookTickerService{c: c, symbol: symbol}
}
func (s *subscribeBookTickerService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@bookTicker"
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsBookTickerResponse struct {
	UpdateID        int64           `json:"u"`
	Symbol          string          `json:"s"`
	BestBidPrice    decimal.Decimal `json:"b"`
	BestBidQuantity decimal.Decimal `json:"B"`
	BestAskPrice    decimal.Decimal `json:"a"`
	BestAskQuantity decimal.Decimal `json:"A"`
}

type subscribeAllBookTickerService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllBookTickerService() *subscribeAllBookTickerService {
	return &subscribeAllBookTickerService{c: c}
}

func (s *subscribeAllBookTickerService) Do(ctx context.Context, handler func(message *WsBookTickerResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + "!bookTicker"
	return request.Subscribe(ctx, s.c, url, handler)
}

type subscribeForceOrderService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeLiquidationOrderService(symbol string) *subscribeForceOrderService {
	return &subscribeForceOrderService{c: c, symbol: symbol}
}

func (s *subscribeForceOrderService) Do(ctx context.Context, handler func(message *WsForceOrderResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@forceOrder"
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsForceOrderResponse struct {
	EventType  string       `json:"e"`
	EventTime  time.Time    `json:"E,format:unixmilli"`
	ForceOrder WsForceOrder `json:"o"`
}

type WsForceOrder struct {
	Symbol                         string          `json:"s"`
	OrderSide                      OrderSide       `json:"S"`
	OrderType                      OrderType       `json:"o"`
	TimeInForce                    TimeInForce     `json:"f"`
	OriginalQuantity               decimal.Decimal `json:"q"`
	Price                          decimal.Decimal `json:"p"`
	AveragePrice                   decimal.Decimal `json:"ap"`
	OrderStatus                    OrderStatus     `json:"X"`
	OrderLastFilledQuantity        decimal.Decimal `json:"l"`
	OrderFilledAccumulatedQuantity decimal.Decimal `json:"z"`
	OrderTradeTime                 time.Time       `json:"T,format:unixmilli"`
}

type subscribeAllForceOrderService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllForceOrderService() *subscribeAllForceOrderService {
	return &subscribeAllForceOrderService{c: c}
}

func (s *subscribeAllForceOrderService) Do(ctx context.Context, handler func(message *WsForceOrderResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + "!forceOrder@arr"
	return request.Subscribe(ctx, s.c, url, handler)
}

type subscribeCombinedDepthService struct {
	c            *FuturesWebSocketClient
	symbolLevels map[string]string
}

// symbolLevels:
// "BTCUSDT": "5@100ms",
// "ETHUSDT": "10@100ms",
func (c *FuturesWebSocketClient) NewSubscribeCombinedDepthService(symbolLevels map[string]string) *subscribeCombinedDepthService {
	return &subscribeCombinedDepthService{c: c, symbolLevels: symbolLevels}
}

func (s *subscribeCombinedDepthService) Do(ctx context.Context, handler func(message *WsCombinedDepthResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
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

type subscribeIncrementalDepthService struct {
	c      *FuturesWebSocketClient
	symbol string
	speed  string
}

// speed: 100ms
func (c *FuturesWebSocketClient) NewSubscribeIncrementalDepthService(symbol string, speed ...string) *subscribeIncrementalDepthService {
	s := subscribeIncrementalDepthService{c: c, symbol: symbol}
	if len(speed) > 0 {
		s.speed = "@" + speed[0]
	}
	return &s
}

func (s *subscribeIncrementalDepthService) Do(ctx context.Context, handler func(message *WsIncrementalDepthResponse, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	url := common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(s.symbol) + "@depth" + s.speed
	return request.Subscribe(ctx, s.c, url, handler)
}

type WsIncrementalDepthResponse struct {
	EventType         string      `json:"e"`
	EventTime         time.Time   `json:"E,format:unixmilli"`
	TransactionTime   time.Time   `json:"T,format:unixmilli"`
	Symbol            string      `json:"s"`
	FirstUpdateID     int64       `json:"U"`
	FinalUpdateID     int64       `json:"u"`
	PrevFinalUpdateID int64       `json:"pu"`
	Bids              []PriceSize `json:"b"`
	Asks              []PriceSize `json:"a"`
}
