package futures

import (
	"context"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// streamPath builds "/ws/<symbol>@<event>". The symbol must be lowercased,
// but event suffixes like "aggTrade" or "bookTicker" are case-sensitive, so
// we only lowercase the part before the first '@'.
func streamPath(stream string) string {
	if i := strings.Index(stream, "@"); i >= 0 {
		return common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(stream[:i]) + stream[i:]
	}
	return common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(stream)
}

// DepthSpeed is an optional update-rate suffix accepted by partial/diff depth
// streams. Use SpeedDefault for the no-suffix (slowest) variant.
type DepthSpeed string

const (
	SpeedDefault DepthSpeed = ""
	Speed500ms   DepthSpeed = "@500ms"
	Speed100ms   DepthSpeed = "@100ms"
)

// SubscribeAggTradeService -- <symbol>@aggTrade
type SubscribeAggTradeService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeAggTradeService(symbol string) *SubscribeAggTradeService {
	return &SubscribeAggTradeService{c: c, symbol: symbol}
}

func (s *SubscribeAggTradeService) Do(ctx context.Context, cb func(*WsAggTradeEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsAggTradeEvent](ctx, s.c, streamPath(s.symbol+"@aggTrade"), cb)
}

type WsAggTradeEvent struct {
	EventType    string          `json:"e"`
	EventTime    time.Time       `json:"E,format:unixmilli"`
	Symbol       string          `json:"s"`
	AggTradeID   int64           `json:"a"`
	Price        decimal.Decimal `json:"p"`
	Quantity     decimal.Decimal `json:"q"`
	FirstTradeID int64           `json:"f"`
	LastTradeID  int64           `json:"l"`
	TradeTime    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}

// SubscribeMarkPriceService -- <symbol>@markPrice (default 3s) or @1s
type SubscribeMarkPriceService struct {
	c      *FuturesWebSocketClient
	symbol string
	fast   bool
}

func (c *FuturesWebSocketClient) NewSubscribeMarkPriceService(symbol string) *SubscribeMarkPriceService {
	return &SubscribeMarkPriceService{c: c, symbol: symbol}
}

// SetFast switches to 1-second updates.
func (s *SubscribeMarkPriceService) SetFast(fast bool) *SubscribeMarkPriceService {
	s.fast = fast
	return s
}

func (s *SubscribeMarkPriceService) Do(ctx context.Context, cb func(*WsMarkPriceEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := s.symbol + "@markPrice"
	if s.fast {
		stream += "@1s"
	}
	return request.Subscribe[WsMarkPriceEvent](ctx, s.c, streamPath(stream), cb)
}

type WsMarkPriceEvent struct {
	EventType            string          `json:"e"`
	EventTime            time.Time       `json:"E,format:unixmilli"`
	Symbol               string          `json:"s"`
	MarkPrice            decimal.Decimal `json:"p"`
	IndexPrice           decimal.Decimal `json:"i"`
	EstimatedSettlePrice decimal.Decimal `json:"P"`
	FundingRate          decimal.Decimal `json:"r"`
	NextFundingTime      time.Time       `json:"T,format:unixmilli"`
}

// SubscribeAllMarkPricesService -- !markPrice@arr (default 3s) or @1s
type SubscribeAllMarkPricesService struct {
	c    *FuturesWebSocketClient
	fast bool
}

func (c *FuturesWebSocketClient) NewSubscribeAllMarkPricesService() *SubscribeAllMarkPricesService {
	return &SubscribeAllMarkPricesService{c: c}
}

func (s *SubscribeAllMarkPricesService) SetFast(fast bool) *SubscribeAllMarkPricesService {
	s.fast = fast
	return s
}

func (s *SubscribeAllMarkPricesService) Do(ctx context.Context, cb func(*[]WsMarkPriceEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := "!markPrice@arr"
	if s.fast {
		stream += "@1s"
	}
	return request.Subscribe[[]WsMarkPriceEvent](ctx, s.c, streamPath(stream), cb)
}

// SubscribeKlineService -- <symbol>@kline_<interval>
type SubscribeKlineService struct {
	c        *FuturesWebSocketClient
	symbol   string
	interval KlineInterval
}

func (c *FuturesWebSocketClient) NewSubscribeKlineService(symbol string, interval KlineInterval) *SubscribeKlineService {
	return &SubscribeKlineService{c: c, symbol: symbol, interval: interval}
}

func (s *SubscribeKlineService) Do(ctx context.Context, cb func(*WsKlineEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsKlineEvent](ctx, s.c, streamPath(s.symbol+"@kline_"+string(s.interval)), cb)
}

type WsKlineEvent struct {
	EventType string    `json:"e"`
	EventTime time.Time `json:"E,format:unixmilli"`
	Symbol    string    `json:"s"`
	Kline     WsKline   `json:"k"`
}

type WsKline struct {
	StartTime                time.Time       `json:"t,format:unixmilli"`
	CloseTime                time.Time       `json:"T,format:unixmilli"`
	Symbol                   string          `json:"s"`
	Interval                 KlineInterval   `json:"i"`
	FirstTradeID             int64           `json:"f"`
	LastTradeID              int64           `json:"L"`
	Open                     decimal.Decimal `json:"o"`
	Close                    decimal.Decimal `json:"c"`
	High                     decimal.Decimal `json:"h"`
	Low                      decimal.Decimal `json:"l"`
	Volume                   decimal.Decimal `json:"v"`
	NumberOfTrades           int64           `json:"n"`
	IsClosed                 bool            `json:"x"`
	QuoteAssetVolume         decimal.Decimal `json:"q"`
	TakerBuyBaseAssetVolume  decimal.Decimal `json:"V"`
	TakerBuyQuoteAssetVolume decimal.Decimal `json:"Q"`
}

// SubscribeMiniTickerService -- <symbol>@miniTicker
type SubscribeMiniTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeMiniTickerService(symbol string) *SubscribeMiniTickerService {
	return &SubscribeMiniTickerService{c: c, symbol: symbol}
}

func (s *SubscribeMiniTickerService) Do(ctx context.Context, cb func(*WsMiniTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsMiniTickerEvent](ctx, s.c, streamPath(s.symbol+"@miniTicker"), cb)
}

type WsMiniTickerEvent struct {
	EventType        string          `json:"e"`
	EventTime        time.Time       `json:"E,format:unixmilli"`
	Symbol           string          `json:"s"`
	ClosePrice       decimal.Decimal `json:"c"`
	OpenPrice        decimal.Decimal `json:"o"`
	HighPrice        decimal.Decimal `json:"h"`
	LowPrice         decimal.Decimal `json:"l"`
	BaseAssetVolume  decimal.Decimal `json:"v"`
	QuoteAssetVolume decimal.Decimal `json:"q"`
}

// SubscribeAllMiniTickersService -- !miniTicker@arr
type SubscribeAllMiniTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllMiniTickersService() *SubscribeAllMiniTickersService {
	return &SubscribeAllMiniTickersService{c: c}
}

func (s *SubscribeAllMiniTickersService) Do(ctx context.Context, cb func(*[]WsMiniTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsMiniTickerEvent](ctx, s.c, streamPath("!miniTicker@arr"), cb)
}

// SubscribeTickerService -- <symbol>@ticker
type SubscribeTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeTickerService(symbol string) *SubscribeTickerService {
	return &SubscribeTickerService{c: c, symbol: symbol}
}

func (s *SubscribeTickerService) Do(ctx context.Context, cb func(*WsTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsTickerEvent](ctx, s.c, streamPath(s.symbol+"@ticker"), cb)
}

type WsTickerEvent struct {
	EventType          string          `json:"e"`
	EventTime          time.Time       `json:"E,format:unixmilli"`
	Symbol             string          `json:"s"`
	PriceChange        decimal.Decimal `json:"p"`
	PriceChangePercent decimal.Decimal `json:"P"`
	WeightedAvgPrice   decimal.Decimal `json:"w"`
	LastPrice          decimal.Decimal `json:"c"`
	LastQty            decimal.Decimal `json:"Q"`
	OpenPrice          decimal.Decimal `json:"o"`
	HighPrice          decimal.Decimal `json:"h"`
	LowPrice           decimal.Decimal `json:"l"`
	BaseAssetVolume    decimal.Decimal `json:"v"`
	QuoteAssetVolume   decimal.Decimal `json:"q"`
	StatsOpenTime      time.Time       `json:"O,format:unixmilli"`
	StatsCloseTime     time.Time       `json:"C,format:unixmilli"`
	FirstTradeID       int64           `json:"F"`
	LastTradeID        int64           `json:"L"`
	TotalTradeCount    int64           `json:"n"`
}

// SubscribeAllTickersService -- !ticker@arr
type SubscribeAllTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllTickersService() *SubscribeAllTickersService {
	return &SubscribeAllTickersService{c: c}
}

func (s *SubscribeAllTickersService) Do(ctx context.Context, cb func(*[]WsTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTickerEvent](ctx, s.c, streamPath("!ticker@arr"), cb)
}

// SubscribeBookTickerService -- <symbol>@bookTicker
type SubscribeBookTickerService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeBookTickerService(symbol string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, symbol: symbol}
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, cb func(*WsBookTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsBookTickerEvent](ctx, s.c, streamPath(s.symbol+"@bookTicker"), cb)
}

type WsBookTickerEvent struct {
	EventType       string          `json:"e"`
	UpdateID        int64           `json:"u"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	TransactionTime time.Time       `json:"T,format:unixmilli"`
	Symbol          string          `json:"s"`
	BidPrice        decimal.Decimal `json:"b"`
	BidQty          decimal.Decimal `json:"B"`
	AskPrice        decimal.Decimal `json:"a"`
	AskQty          decimal.Decimal `json:"A"`
}

// SubscribeAllBookTickersService -- !bookTicker
type SubscribeAllBookTickersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllBookTickersService() *SubscribeAllBookTickersService {
	return &SubscribeAllBookTickersService{c: c}
}

func (s *SubscribeAllBookTickersService) Do(ctx context.Context, cb func(*WsBookTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsBookTickerEvent](ctx, s.c, streamPath("!bookTicker"), cb)
}

// SubscribeForceOrderService -- <symbol>@forceOrder
type SubscribeForceOrderService struct {
	c      *FuturesWebSocketClient
	symbol string
}

func (c *FuturesWebSocketClient) NewSubscribeForceOrderService(symbol string) *SubscribeForceOrderService {
	return &SubscribeForceOrderService{c: c, symbol: symbol}
}

func (s *SubscribeForceOrderService) Do(ctx context.Context, cb func(*WsForceOrderEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsForceOrderEvent](ctx, s.c, streamPath(s.symbol+"@forceOrder"), cb)
}

type WsForceOrderEvent struct {
	EventType string         `json:"e"`
	EventTime time.Time      `json:"E,format:unixmilli"`
	Order     WsForceOrderRO `json:"o"`
}

type WsForceOrderRO struct {
	Symbol            string          `json:"s"`
	Side              OrderSide       `json:"S"`
	OrderType         OrderType       `json:"o"`
	TimeInForce       TimeInForce     `json:"f"`
	OrigQty           decimal.Decimal `json:"q"`
	Price             decimal.Decimal `json:"p"`
	AvgPrice          decimal.Decimal `json:"ap"`
	Status            OrderStatus     `json:"X"`
	LastFilledQty     decimal.Decimal `json:"l"`
	AccumFilledQty    decimal.Decimal `json:"z"`
	OrderTradeTime    time.Time       `json:"T,format:unixmilli"`
}

// SubscribeAllForceOrdersService -- !forceOrder@arr
type SubscribeAllForceOrdersService struct {
	c *FuturesWebSocketClient
}

func (c *FuturesWebSocketClient) NewSubscribeAllForceOrdersService() *SubscribeAllForceOrdersService {
	return &SubscribeAllForceOrdersService{c: c}
}

func (s *SubscribeAllForceOrdersService) Do(ctx context.Context, cb func(*WsForceOrderEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsForceOrderEvent](ctx, s.c, streamPath("!forceOrder@arr"), cb)
}

// SubscribePartialDepthService -- <symbol>@depth<levels>[@500ms|@100ms]
type SubscribePartialDepthService struct {
	c      *FuturesWebSocketClient
	symbol string
	levels int
	speed  DepthSpeed
}

func (c *FuturesWebSocketClient) NewSubscribePartialDepthService(symbol string, levels int) *SubscribePartialDepthService {
	return &SubscribePartialDepthService{c: c, symbol: symbol, levels: levels}
}

func (s *SubscribePartialDepthService) SetSpeed(speed DepthSpeed) *SubscribePartialDepthService {
	s.speed = speed
	return s
}

func (s *SubscribePartialDepthService) Do(ctx context.Context, cb func(*WsDepthEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := s.symbol + "@depth" + intToStr(s.levels) + string(s.speed)
	return request.Subscribe[WsDepthEvent](ctx, s.c, streamPath(stream), cb)
}

// SubscribeDiffDepthService -- <symbol>@depth[@500ms|@100ms]
//
// Diff streams use the same payload shape as partial depth on V3 futures
// (b/a lowercase), unlike the V3 spot diff/partial split.
type SubscribeDiffDepthService struct {
	c      *FuturesWebSocketClient
	symbol string
	speed  DepthSpeed
}

func (c *FuturesWebSocketClient) NewSubscribeDiffDepthService(symbol string) *SubscribeDiffDepthService {
	return &SubscribeDiffDepthService{c: c, symbol: symbol}
}

func (s *SubscribeDiffDepthService) SetSpeed(speed DepthSpeed) *SubscribeDiffDepthService {
	s.speed = speed
	return s
}

func (s *SubscribeDiffDepthService) Do(ctx context.Context, cb func(*WsDepthEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := s.symbol + "@depth" + string(s.speed)
	return request.Subscribe[WsDepthEvent](ctx, s.c, streamPath(stream), cb)
}

type WsDepthEvent struct {
	EventType       string      `json:"e"`
	EventTime       time.Time   `json:"E,format:unixmilli"`
	TransactionTime time.Time   `json:"T,format:unixmilli"`
	Symbol          string      `json:"s"`
	FirstUpdateID   int64       `json:"U"`
	FinalUpdateID   int64       `json:"u"`
	PrevFinalID     int64       `json:"pu"`
	Bids            [][2]string `json:"b"`
	Asks            [][2]string `json:"a"`
}

// intToStr is a tiny helper avoiding the strconv import for small ints used in
// stream names. Levels are typically 5/10/20.
func intToStr(i int) string {
	switch i {
	case 0:
		return "0"
	case 5:
		return "5"
	case 10:
		return "10"
	case 20:
		return "20"
	}
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	n := len(buf)
	for i > 0 {
		n--
		buf[n] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		n--
		buf[n] = '-'
	}
	return string(buf[n:])
}
