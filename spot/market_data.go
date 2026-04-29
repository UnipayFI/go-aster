package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// GetExchangeInfoService -- GET /api/v3/exchangeInfo
type GetExchangeInfoService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetExchangeInfoService() *GetExchangeInfoService {
	return &GetExchangeInfoService{c: c}
}

func (s *GetExchangeInfoService) Do(ctx context.Context) (*ExchangeInfoResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/exchangeInfo")
	return request.Do[ExchangeInfoResponse](req)
}

type ExchangeInfoResponse struct {
	Timezone        string      `json:"timezone"`
	ServerTime      time.Time   `json:"serverTime,format:unixmilli"`
	RateLimits      []RateLimit `json:"rateLimits"`
	ExchangeFilters []any       `json:"exchangeFilters"`
	Assets          []AssetInfo `json:"assets"`
	Symbols         []Symbol    `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

type AssetInfo struct {
	Asset string `json:"asset"`
}

type Symbol struct {
	Symbol             string           `json:"symbol"`
	Status             string           `json:"status"`
	BaseAsset          string           `json:"baseAsset"`
	QuoteAsset         string           `json:"quoteAsset"`
	PricePrecision     int              `json:"pricePrecision"`
	QuantityPrecision  int              `json:"quantityPrecision"`
	BaseAssetPrecision int              `json:"baseAssetPrecision"`
	QuotePrecision     int              `json:"quotePrecision"`
	Filters            []map[string]any `json:"filters"`
	OrderTypes         []OrderType      `json:"orderTypes"`
	TimeInForce        []TimeInForce    `json:"timeInForce"`
	OcoAllowed         bool             `json:"ocoAllowed"`
}

// GetDepthService -- GET /api/v3/depth
type GetDepthService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetDepthService(symbol string) *GetDepthService {
	return &GetDepthService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetDepthService) SetLimit(limit int) *GetDepthService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetDepthService) Do(ctx context.Context) (*DepthResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/depth", s.params)
	return request.Do[DepthResponse](req)
}

// DepthResponse -- price levels are [price, qty] string pairs to preserve full
// precision; convert with shopspring/decimal at the call site if needed.
type DepthResponse struct {
	LastUpdateId    int64       `json:"lastUpdateId"`
	MessageTime     time.Time   `json:"E,format:unixmilli"`
	TransactionTime time.Time   `json:"T,format:unixmilli"`
	Bids            [][2]string `json:"bids"`
	Asks            [][2]string `json:"asks"`
}

// GetRecentTradesService -- GET /api/v3/trades
type GetRecentTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetRecentTradesService(symbol string) *GetRecentTradesService {
	return &GetRecentTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetRecentTradesService) SetLimit(limit int) *GetRecentTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecentTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/api/v3/trades", s.params)
	resp, err := request.Do[[]Trade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type Trade struct {
	ID           int64           `json:"id"`
	Price        decimal.Decimal `json:"price"`
	Qty          decimal.Decimal `json:"qty"`
	BaseQty      decimal.Decimal `json:"baseQty"`
	Time         time.Time       `json:"time,format:unixmilli"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
}

// GetHistoricalTradesService -- GET /api/v3/historicalTrades (MARKET_DATA)
type GetHistoricalTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetHistoricalTradesService(symbol string) *GetHistoricalTradesService {
	return &GetHistoricalTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetHistoricalTradesService) SetLimit(limit int) *GetHistoricalTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoricalTradesService) SetFromId(fromId int64) *GetHistoricalTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *GetHistoricalTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/api/v3/historicalTrades", s.params)
	resp, err := request.Do[[]Trade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetAggTradesService -- GET /api/v3/aggTrades
type GetAggTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAggTradesService(symbol string) *GetAggTradesService {
	return &GetAggTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetAggTradesService) SetFromId(fromId int64) *GetAggTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *GetAggTradesService) SetStartTime(t time.Time) *GetAggTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAggTradesService) SetEndTime(t time.Time) *GetAggTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAggTradesService) SetLimit(limit int) *GetAggTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAggTradesService) Do(ctx context.Context) ([]AggTrade, error) {
	req := request.Get(ctx, s.c, "/api/v3/aggTrades", s.params)
	resp, err := request.Do[[]AggTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type AggTrade struct {
	AggTradeID   int64           `json:"a"`
	Price        decimal.Decimal `json:"p"`
	Qty          decimal.Decimal `json:"q"`
	FirstTradeID int64           `json:"f"`
	LastTradeID  int64           `json:"l"`
	Timestamp    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}

// GetKlinesService -- GET /api/v3/klines
type GetKlinesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetKlinesService(symbol string, interval KlineInterval) *GetKlinesService {
	return &GetKlinesService{c: c, params: map[string]string{
		"symbol":   symbol,
		"interval": string(interval),
	}}
}

func (s *GetKlinesService) SetStartTime(t time.Time) *GetKlinesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetKlinesService) SetEndTime(t time.Time) *GetKlinesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetKlinesService) SetLimit(limit int) *GetKlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// Kline format: API returns a fixed-shape array per candle. We keep the raw
// `[]any` form and provide a typed view via parseKline; this avoids losing
// precision on the decimal fields when go's json experiments coerce numbers.
func (s *GetKlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, "/api/v3/klines", s.params)
	raw, err := request.Do[[][]any](req)
	if err != nil {
		return nil, err
	}
	out := make([]Kline, 0, len(*raw))
	for _, row := range *raw {
		k, err := parseKline(row)
		if err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, nil
}

type Kline struct {
	OpenTime                 time.Time
	Open                     decimal.Decimal
	High                     decimal.Decimal
	Low                      decimal.Decimal
	Close                    decimal.Decimal
	Volume                   decimal.Decimal
	CloseTime                time.Time
	QuoteAssetVolume         decimal.Decimal
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  decimal.Decimal
	TakerBuyQuoteAssetVolume decimal.Decimal
}

func parseKline(row []any) (Kline, error) {
	asDecimal := func(v any) (decimal.Decimal, error) {
		s, ok := v.(string)
		if !ok {
			return decimal.Zero, errInvalidKlineRow
		}
		return decimal.NewFromString(s)
	}
	asInt := func(v any) (int64, error) {
		switch x := v.(type) {
		case float64:
			return int64(x), nil
		case int64:
			return x, nil
		case int:
			return int64(x), nil
		}
		return 0, errInvalidKlineRow
	}
	if len(row) < 11 {
		return Kline{}, errInvalidKlineRow
	}
	openMs, err := asInt(row[0])
	if err != nil {
		return Kline{}, err
	}
	closeMs, err := asInt(row[6])
	if err != nil {
		return Kline{}, err
	}
	open, err := asDecimal(row[1])
	if err != nil {
		return Kline{}, err
	}
	high, err := asDecimal(row[2])
	if err != nil {
		return Kline{}, err
	}
	low, err := asDecimal(row[3])
	if err != nil {
		return Kline{}, err
	}
	cl, err := asDecimal(row[4])
	if err != nil {
		return Kline{}, err
	}
	vol, err := asDecimal(row[5])
	if err != nil {
		return Kline{}, err
	}
	quoteVol, err := asDecimal(row[7])
	if err != nil {
		return Kline{}, err
	}
	numTrades, err := asInt(row[8])
	if err != nil {
		return Kline{}, err
	}
	takerBase, err := asDecimal(row[9])
	if err != nil {
		return Kline{}, err
	}
	takerQuote, err := asDecimal(row[10])
	if err != nil {
		return Kline{}, err
	}
	return Kline{
		OpenTime:                 time.UnixMilli(openMs),
		Open:                     open,
		High:                     high,
		Low:                      low,
		Close:                    cl,
		Volume:                   vol,
		CloseTime:                time.UnixMilli(closeMs),
		QuoteAssetVolume:         quoteVol,
		NumberOfTrades:           numTrades,
		TakerBuyBaseAssetVolume:  takerBase,
		TakerBuyQuoteAssetVolume: takerQuote,
	}, nil
}

var errInvalidKlineRow = errInvalidRow("invalid kline row")

type errInvalidRow string

func (e errInvalidRow) Error() string { return string(e) }

// Get24hTickerService -- GET /api/v3/ticker/24hr
//
// When symbol is empty the API returns an array; when set it returns a single
// object. To keep the API uniform, the typed accessor returns []Ticker24h.
type Get24hTickerService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGet24hTickerService() *Get24hTickerService {
	return &Get24hTickerService{c: c, params: map[string]string{}}
}

func (s *Get24hTickerService) SetSymbol(symbol string) *Get24hTickerService {
	s.params["symbol"] = symbol
	return s
}

func (s *Get24hTickerService) Do(ctx context.Context) ([]Ticker24h, error) {
	req := request.Get(ctx, s.c, "/api/v3/ticker/24hr", s.params)
	if _, ok := s.params["symbol"]; ok {
		single, err := request.Do[Ticker24h](req)
		if err != nil {
			return nil, err
		}
		return []Ticker24h{*single}, nil
	}
	multi, err := request.Do[[]Ticker24h](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type Ticker24h struct {
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
	OpenTime           time.Time       `json:"openTime,format:unixmilli"`
	CloseTime          time.Time       `json:"closeTime,format:unixmilli"`
	FirstId            int64           `json:"firstId"`
	LastId             int64           `json:"lastId"`
	Count              int64           `json:"count"`
	BaseAsset          string          `json:"baseAsset"`
	QuoteAsset         string          `json:"quoteAsset"`
}

// GetTickerPriceService -- GET /api/v3/ticker/price
type GetTickerPriceService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTickerPriceService() *GetTickerPriceService {
	return &GetTickerPriceService{c: c, params: map[string]string{}}
}

func (s *GetTickerPriceService) SetSymbol(symbol string) *GetTickerPriceService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetTickerPriceService) Do(ctx context.Context) ([]TickerPrice, error) {
	req := request.Get(ctx, s.c, "/api/v3/ticker/price", s.params)
	if _, ok := s.params["symbol"]; ok {
		single, err := request.Do[TickerPrice](req)
		if err != nil {
			return nil, err
		}
		return []TickerPrice{*single}, nil
	}
	multi, err := request.Do[[]TickerPrice](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type TickerPrice struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Time   time.Time       `json:"time,format:unixmilli"`
}

// GetBookTickerService -- GET /api/v3/ticker/bookTicker
type GetBookTickerService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetBookTickerService() *GetBookTickerService {
	return &GetBookTickerService{c: c, params: map[string]string{}}
}

func (s *GetBookTickerService) SetSymbol(symbol string) *GetBookTickerService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetBookTickerService) Do(ctx context.Context) ([]BookTicker, error) {
	req := request.Get(ctx, s.c, "/api/v3/ticker/bookTicker", s.params)
	if _, ok := s.params["symbol"]; ok {
		single, err := request.Do[BookTicker](req)
		if err != nil {
			return nil, err
		}
		return []BookTicker{*single}, nil
	}
	multi, err := request.Do[[]BookTicker](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type BookTicker struct {
	Symbol   string          `json:"symbol"`
	BidPrice decimal.Decimal `json:"bidPrice"`
	BidQty   decimal.Decimal `json:"bidQty"`
	AskPrice decimal.Decimal `json:"askPrice"`
	AskQty   decimal.Decimal `json:"askQty"`
	Time     time.Time       `json:"time,format:unixmilli"`
}

// GetCommissionRateService -- GET /api/v3/commissionRate (USER_DATA)
type GetCommissionRateService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetCommissionRateService(symbol string) *GetCommissionRateService {
	return &GetCommissionRateService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetCommissionRateService) Do(ctx context.Context) (*CommissionRateResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/commissionRate", s.params).WithSignature()
	return request.Do[CommissionRateResponse](req)
}

type CommissionRateResponse struct {
	Symbol              string          `json:"symbol"`
	MakerCommissionRate decimal.Decimal `json:"makerCommissionRate"`
	TakerCommissionRate decimal.Decimal `json:"takerCommissionRate"`
}
