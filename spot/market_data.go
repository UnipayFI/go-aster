package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type DepthService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewDepthService(symbol string) *DepthService {
	return &DepthService{c: c, params: map[string]string{"symbol": symbol}}
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
	req := request.Get(ctx, s.c, "/api/v1/depth", s.params)
	return request.Do[DepthResponse](req)
}

type DepthResponse struct {
	LastUpdateId int64                `json:"lastUpdateId"`
	Bids         [][2]decimal.Decimal `json:"bids"`
	Asks         [][2]decimal.Decimal `json:"asks"`
}

type GetTradesService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewGetTradesService(symbol string) *GetTradesService {
	return &GetTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetTradesService) SetLimit(limit int) *GetTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetTradesService) Do(ctx context.Context) ([]TradeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/trades", s.params)
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

type HistoricalTradesService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewHistoricalTradesService(symbol string) *HistoricalTradesService {
	return &HistoricalTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *HistoricalTradesService) SetLimit(limit int) *HistoricalTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *HistoricalTradesService) SetFromId(fromId int64) *HistoricalTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *HistoricalTradesService) Do(ctx context.Context) ([]TradeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/historicalTrades", s.params).SetApiKeyHeader()
	trades, err := request.Do[[]TradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type AggTradesService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewAggTradesService(symbol string) *AggTradesService {
	return &AggTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *AggTradesService) SetFromId(fromId int64) *AggTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *AggTradesService) SetStartTime(startTime int64) *AggTradesService {
	s.params["startTime"] = strconv.FormatInt(startTime, 10)
	return s
}

func (s *AggTradesService) SetEndTime(endTime int64) *AggTradesService {
	s.params["endTime"] = strconv.FormatInt(endTime, 10)
	return s
}

func (s *AggTradesService) SetLimit(limit int) *AggTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *AggTradesService) Do(ctx context.Context) ([]AggTradeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/aggTrades", s.params)
	trades, err := request.Do[[]AggTradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type AggTradeResponse struct {
	AggTradeId   int64           `json:"a"`
	Price        decimal.Decimal `json:"p"`
	Quantity     decimal.Decimal `json:"q"`
	FirstTradeId int64           `json:"f"`
	LastTradeId  int64           `json:"l"`
	Timestamp    int64           `json:"T"`
	IsBuyerMaker bool            `json:"m"`
}

type KlinesService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewKlinesService(symbol string, interval KlineInterval) *KlinesService {
	return &KlinesService{c: c, params: map[string]string{"symbol": symbol, "interval": string(interval)}}
}

func (s *KlinesService) SetStartTime(startTime int64) *KlinesService {
	s.params["startTime"] = strconv.FormatInt(startTime, 10)
	return s
}

func (s *KlinesService) SetEndTime(endTime int64) *KlinesService {
	s.params["endTime"] = strconv.FormatInt(endTime, 10)
	return s
}

func (s *KlinesService) SetLimit(limit int) *KlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *KlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, "/api/v1/klines", s.params)
	klines, err := request.Do[[]Kline](req)
	if err != nil {
		return nil, err
	}
	return *klines, nil
}

type Kline = []any

type Ticker24hrService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewTicker24hrService() *Ticker24hrService {
	return &Ticker24hrService{c: c, params: map[string]string{}}
}

func (s *Ticker24hrService) SetSymbol(symbol string) *Ticker24hrService {
	s.params["symbol"] = symbol
	return s
}

func (s *Ticker24hrService) Do(ctx context.Context) (*Ticker24hrResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/24hr", s.params)
	return request.Do[Ticker24hrResponse](req)
}

func (s *Ticker24hrService) DoAll(ctx context.Context) ([]Ticker24hrResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/24hr", s.params)
	tickers, err := request.Do[[]Ticker24hrResponse](req)
	if err != nil {
		return nil, err
	}
	return *tickers, nil
}

type Ticker24hrResponse struct {
	Symbol             string          `json:"symbol"`
	PriceChange        decimal.Decimal `json:"priceChange"`
	PriceChangePercent decimal.Decimal `json:"priceChangePercent"`
	WeightedAvgPrice   decimal.Decimal `json:"weightedAvgPrice"`
	PrevClosePrice     decimal.Decimal `json:"prevClosePrice"`
	LastPrice          decimal.Decimal `json:"lastPrice"`
	LastQty            decimal.Decimal `json:"lastQty"`
	BidPrice           decimal.Decimal `json:"bidPrice"`
	BidQty             decimal.Decimal `json:"bidQty"`
	AskPrice           decimal.Decimal `json:"askPrice"`
	AskQty             decimal.Decimal `json:"askQty"`
	OpenPrice          decimal.Decimal `json:"openPrice"`
	HighPrice          decimal.Decimal `json:"highPrice"`
	LowPrice           decimal.Decimal `json:"lowPrice"`
	Volume             decimal.Decimal `json:"volume"`
	QuoteVolume        decimal.Decimal `json:"quoteVolume"`
	OpenTime           int64           `json:"openTime"`
	CloseTime          int64           `json:"closeTime"`
	FirstId            int64           `json:"firstId"`
	LastId             int64           `json:"lastId"`
	Count              int64           `json:"count"`
	BaseAsset          string          `json:"baseAsset"`
	QuoteAsset         string          `json:"quoteAsset"`
}

type TickerPriceService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewTickerPriceService() *TickerPriceService {
	return &TickerPriceService{c: c, params: map[string]string{}}
}

func (s *TickerPriceService) SetSymbol(symbol string) *TickerPriceService {
	s.params["symbol"] = symbol
	return s
}

func (s *TickerPriceService) Do(ctx context.Context) (*TickerPriceResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/price", s.params)
	return request.Do[TickerPriceResponse](req)
}

func (s *TickerPriceService) DoAll(ctx context.Context) ([]TickerPriceResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/price", s.params)
	tickers, err := request.Do[[]TickerPriceResponse](req)
	if err != nil {
		return nil, err
	}
	return *tickers, nil
}

type TickerPriceResponse struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Time   int64           `json:"time"`
}

type BookTickerService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewBookTickerService() *BookTickerService {
	return &BookTickerService{c: c, params: map[string]string{}}
}

func (s *BookTickerService) SetSymbol(symbol string) *BookTickerService {
	s.params["symbol"] = symbol
	return s
}

func (s *BookTickerService) Do(ctx context.Context) (*BookTickerResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/bookTicker", s.params)
	return request.Do[BookTickerResponse](req)
}

func (s *BookTickerService) DoAll(ctx context.Context) ([]BookTickerResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/ticker/bookTicker", s.params)
	tickers, err := request.Do[[]BookTickerResponse](req)
	if err != nil {
		return nil, err
	}
	return *tickers, nil
}

type BookTickerResponse struct {
	Symbol   string          `json:"symbol"`
	BidPrice decimal.Decimal `json:"bidPrice"`
	BidQty   decimal.Decimal `json:"bidQty"`
	AskPrice decimal.Decimal `json:"askPrice"`
	AskQty   decimal.Decimal `json:"askQty"`
	Time     int64           `json:"time"`
}

type CommissionRateService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewCommissionRateService(symbol string) *CommissionRateService {
	return &CommissionRateService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *CommissionRateService) Do(ctx context.Context) (*CommissionRateResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/commissionRate", s.params).Sign()
	return request.Do[CommissionRateResponse](req)
}

type CommissionRateResponse struct {
	Symbol              string          `json:"symbol"`
	MakerCommissionRate decimal.Decimal `json:"makerCommissionRate"`
	TakerCommissionRate decimal.Decimal `json:"takerCommissionRate"`
}
