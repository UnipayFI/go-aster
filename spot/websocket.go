package spot

import (
	"context"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// streamPath builds a single-stream URL ("/ws/<symbol>@<event>"). Aster
// requires lowercase symbols, but event suffixes like "aggTrade" or
// "bookTicker" are case-sensitive — only the segment before the first '@'
// gets lowercased.
func streamPath(stream string) string {
	if i := strings.Index(stream, "@"); i >= 0 {
		return common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(stream[:i]) + stream[i:]
	}
	return common.WEBSOCKET_STREAM_SEPARATOR + strings.ToLower(stream)
}

// SubscribeAggTradeService -- <symbol>@aggTrade (real-time aggregated trades)
type SubscribeAggTradeService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeAggTradeService(symbol string) *SubscribeAggTradeService {
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

// SubscribeTradeService -- <symbol>@trade (tick-by-tick raw trades)
type SubscribeTradeService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeTradeService(symbol string) *SubscribeTradeService {
	return &SubscribeTradeService{c: c, symbol: symbol}
}

func (s *SubscribeTradeService) Do(ctx context.Context, cb func(*WsTradeEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsTradeEvent](ctx, s.c, streamPath(s.symbol+"@trade"), cb)
}

type WsTradeEvent struct {
	EventType    string          `json:"e"`
	EventTime    time.Time       `json:"E,format:unixmilli"`
	Symbol       string          `json:"s"`
	TradeID      int64           `json:"t"`
	Price        decimal.Decimal `json:"p"`
	Quantity     decimal.Decimal `json:"q"`
	TradeTime    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}

// SubscribeKlineService -- <symbol>@kline_<interval>
type SubscribeKlineService struct {
	c        *SpotWebSocketClient
	symbol   string
	interval KlineInterval
}

func (c *SpotWebSocketClient) NewSubscribeKlineService(symbol string, interval KlineInterval) *SubscribeKlineService {
	return &SubscribeKlineService{c: c, symbol: symbol, interval: interval}
}

func (s *SubscribeKlineService) Do(ctx context.Context, cb func(*WsKlineEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsKlineEvent](ctx, s.c, streamPath(s.symbol+"@kline_"+string(s.interval)), cb)
}

type WsKlineEvent struct {
	EventType string  `json:"e"`
	EventTime time.Time `json:"E,format:unixmilli"`
	Symbol    string  `json:"s"`
	Kline     WsKline `json:"k"`
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
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeMiniTickerService(symbol string) *SubscribeMiniTickerService {
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

// SubscribeAllMiniTickersService -- !miniTicker@arr (all symbols)
type SubscribeAllMiniTickersService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllMiniTickersService() *SubscribeAllMiniTickersService {
	return &SubscribeAllMiniTickersService{c: c}
}

func (s *SubscribeAllMiniTickersService) Do(ctx context.Context, cb func(*[]WsMiniTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsMiniTickerEvent](ctx, s.c, streamPath("!miniTicker@arr"), cb)
}

// SubscribeTickerService -- <symbol>@ticker (full 24h ticker)
type SubscribeTickerService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeTickerService(symbol string) *SubscribeTickerService {
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
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllTickersService() *SubscribeAllTickersService {
	return &SubscribeAllTickersService{c: c}
}

func (s *SubscribeAllTickersService) Do(ctx context.Context, cb func(*[]WsTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTickerEvent](ctx, s.c, streamPath("!ticker@arr"), cb)
}

// SubscribeBookTickerService -- <symbol>@bookTicker
type SubscribeBookTickerService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeBookTickerService(symbol string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, symbol: symbol}
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, cb func(*WsBookTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsBookTickerEvent](ctx, s.c, streamPath(s.symbol+"@bookTicker"), cb)
}

type WsBookTickerEvent struct {
	UpdateID int64           `json:"u"`
	Symbol   string          `json:"s"`
	BidPrice decimal.Decimal `json:"b"`
	BidQty   decimal.Decimal `json:"B"`
	AskPrice decimal.Decimal `json:"a"`
	AskQty   decimal.Decimal `json:"A"`
}

// SubscribeAllBookTickersService -- !bookTicker
type SubscribeAllBookTickersService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeAllBookTickersService() *SubscribeAllBookTickersService {
	return &SubscribeAllBookTickersService{c: c}
}

func (s *SubscribeAllBookTickersService) Do(ctx context.Context, cb func(*WsBookTickerEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsBookTickerEvent](ctx, s.c, streamPath("!bookTicker"), cb)
}

// SubscribePartialDepthService -- <symbol>@depth<levels> or <symbol>@depth<levels>@100ms
//
// levels must be 5, 10, or 20. Set fast=true for 100ms updates.
type SubscribePartialDepthService struct {
	c      *SpotWebSocketClient
	symbol string
	levels int
	fast   bool
}

func (c *SpotWebSocketClient) NewSubscribePartialDepthService(symbol string, levels int) *SubscribePartialDepthService {
	return &SubscribePartialDepthService{c: c, symbol: symbol, levels: levels}
}

func (s *SubscribePartialDepthService) SetFast(fast bool) *SubscribePartialDepthService {
	s.fast = fast
	return s
}

func (s *SubscribePartialDepthService) Do(ctx context.Context, cb func(*WsPartialDepthEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := s.symbol + "@depth" + itoa(s.levels)
	if s.fast {
		stream += "@100ms"
	}
	return request.Subscribe[WsPartialDepthEvent](ctx, s.c, streamPath(stream), cb)
}

type WsPartialDepthEvent struct {
	EventType       string      `json:"e"`
	EventTime       time.Time   `json:"E,format:unixmilli"`
	TransactionTime time.Time   `json:"T,format:unixmilli"`
	Symbol          string      `json:"s"`
	FirstUpdateID   int64       `json:"U"`
	FinalUpdateID   int64       `json:"u"`
	PrevFinalID     int64       `json:"pu"`
	Bids            [][2]string `json:"bids"`
	Asks            [][2]string `json:"asks"`
}

// SubscribeDiffDepthService -- <symbol>@depth or <symbol>@depth@100ms
//
// Note the payload uses "b"/"a" (lowercase short form) for diff updates,
// distinct from PartialDepth's "bids"/"asks".
type SubscribeDiffDepthService struct {
	c      *SpotWebSocketClient
	symbol string
	fast   bool
}

func (c *SpotWebSocketClient) NewSubscribeDiffDepthService(symbol string) *SubscribeDiffDepthService {
	return &SubscribeDiffDepthService{c: c, symbol: symbol}
}

func (s *SubscribeDiffDepthService) SetFast(fast bool) *SubscribeDiffDepthService {
	s.fast = fast
	return s
}

func (s *SubscribeDiffDepthService) Do(ctx context.Context, cb func(*WsDiffDepthEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	stream := s.symbol + "@depth"
	if s.fast {
		stream += "@100ms"
	}
	return request.Subscribe[WsDiffDepthEvent](ctx, s.c, streamPath(stream), cb)
}

type WsDiffDepthEvent struct {
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

// SubscribeTradeProService -- <symbol>@tradepro
//
// Aster-specific stream that exposes the on-chain transaction hash and the
// taker/maker addresses for each fill (or "hidden" for hidden orders).
type SubscribeTradeProService struct {
	c      *SpotWebSocketClient
	symbol string
}

func (c *SpotWebSocketClient) NewSubscribeTradeProService(symbol string) *SubscribeTradeProService {
	return &SubscribeTradeProService{c: c, symbol: symbol}
}

func (s *SubscribeTradeProService) Do(ctx context.Context, cb func(*WsTradeProEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsTradeProEvent](ctx, s.c, streamPath(s.symbol+"@tradepro"), cb)
}

type WsTradeProEvent struct {
	EventType       string          `json:"e"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	TradeTime       time.Time       `json:"T,format:unixmilli"`
	Symbol          string          `json:"s"`
	TradeID         int64           `json:"t"`
	Price           decimal.Decimal `json:"p"`
	Quantity        decimal.Decimal `json:"q"`
	TransactionHash string          `json:"h"`
	Participants    []string        `json:"m"` // [taker, maker], "hidden" if order was hidden
}

func itoa(i int) string {
	// small helper to avoid importing strconv just for one digit
	switch i {
	case 5:
		return "5"
	case 10:
		return "10"
	case 20:
		return "20"
	}
	// fall back to a sane format for unexpected levels
	return integerToString(i)
}

func integerToString(i int) string {
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
