package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// GetExchangeInfoService -- GET /fapi/v3/exchangeInfo
type GetExchangeInfoService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetExchangeInfoService() *GetExchangeInfoService {
	return &GetExchangeInfoService{c: c}
}

func (s *GetExchangeInfoService) Do(ctx context.Context) (*ExchangeInfoResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/exchangeInfo")
	return request.Do[ExchangeInfoResponse](req)
}

type ExchangeInfoResponse struct {
	Timezone        string          `json:"timezone"`
	ServerTime      time.Time       `json:"serverTime,format:unixmilli"`
	RateLimits      []RateLimit     `json:"rateLimits"`
	ExchangeFilters []any           `json:"exchangeFilters"`
	Assets          []AssetInfo     `json:"assets"`
	Symbols         []FuturesSymbol `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

type AssetInfo struct {
	Asset             string          `json:"asset"`
	MarginAvailable   bool            `json:"marginAvailable"`
	AutoAssetExchange decimal.Decimal `json:"autoAssetExchange"`
}

type FuturesSymbol struct {
	Symbol                string           `json:"symbol"`
	Pair                  string           `json:"pair"`
	ContractType          string           `json:"contractType"`
	DeliveryDate          int64            `json:"deliveryDate"`
	OnboardDate           int64            `json:"onboardDate"`
	Status                string           `json:"status"`
	MaintMarginPercent    decimal.Decimal  `json:"maintMarginPercent"`
	RequiredMarginPercent decimal.Decimal  `json:"requiredMarginPercent"`
	BaseAsset             string           `json:"baseAsset"`
	QuoteAsset            string           `json:"quoteAsset"`
	MarginAsset           string           `json:"marginAsset"`
	PricePrecision        int              `json:"pricePrecision"`
	QuantityPrecision     int              `json:"quantityPrecision"`
	BaseAssetPrecision    int              `json:"baseAssetPrecision"`
	QuotePrecision        int              `json:"quotePrecision"`
	UnderlyingType        string           `json:"underlyingType"`
	UnderlyingSubType     []string         `json:"underlyingSubType"`
	SettlePlan            int              `json:"settlePlan"`
	TriggerProtect        decimal.Decimal  `json:"triggerProtect"`
	Filters               []map[string]any `json:"filters"`
	OrderTypes            []OrderType      `json:"OrderType"`
	TimeInForce           []TimeInForce    `json:"timeInForce"`
	LiquidationFee        decimal.Decimal  `json:"liquidationFee"`
	MarketTakeBound       decimal.Decimal  `json:"marketTakeBound"`
}

// GetDepthService -- GET /fapi/v3/depth
type GetDepthService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetDepthService(symbol string) *GetDepthService {
	return &GetDepthService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetDepthService) SetLimit(limit int) *GetDepthService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetDepthService) Do(ctx context.Context) (*DepthResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/depth", s.params)
	return request.Do[DepthResponse](req)
}

type DepthResponse struct {
	LastUpdateId    int64       `json:"lastUpdateId"`
	MessageTime     time.Time   `json:"E,format:unixmilli"`
	TransactionTime time.Time   `json:"T,format:unixmilli"`
	Bids            [][2]string `json:"bids"`
	Asks            [][2]string `json:"asks"`
}

// GetRecentTradesService -- GET /fapi/v3/trades
type GetRecentTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetRecentTradesService(symbol string) *GetRecentTradesService {
	return &GetRecentTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetRecentTradesService) SetLimit(limit int) *GetRecentTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecentTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/trades", s.params)
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
	QuoteQty     decimal.Decimal `json:"quoteQty"`
	Time         time.Time       `json:"time,format:unixmilli"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
}

// GetHistoricalTradesService -- GET /fapi/v3/historicalTrades
type GetHistoricalTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetHistoricalTradesService(symbol string) *GetHistoricalTradesService {
	return &GetHistoricalTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetHistoricalTradesService) SetLimit(limit int) *GetHistoricalTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoricalTradesService) SetFromId(id int64) *GetHistoricalTradesService {
	s.params["fromId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetHistoricalTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/historicalTrades", s.params)
	resp, err := request.Do[[]Trade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetAggTradesService -- GET /fapi/v3/aggTrades
type GetAggTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetAggTradesService(symbol string) *GetAggTradesService {
	return &GetAggTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetAggTradesService) SetFromId(id int64) *GetAggTradesService {
	s.params["fromId"] = strconv.FormatInt(id, 10)
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
	req := request.Get(ctx, s.c, "/fapi/v3/aggTrades", s.params)
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

// GetKlinesService -- GET /fapi/v3/klines
type GetKlinesService struct {
	c      *FuturesClient
	params map[string]string
	path   string
}

func (c *FuturesClient) NewGetKlinesService(symbol string, interval KlineInterval) *GetKlinesService {
	return &GetKlinesService{c: c, path: "/fapi/v3/klines", params: map[string]string{
		"symbol":   symbol,
		"interval": string(interval),
	}}
}

// NewGetIndexPriceKlinesService -- GET /fapi/v3/indexPriceKlines (note: pair, not symbol)
func (c *FuturesClient) NewGetIndexPriceKlinesService(pair string, interval KlineInterval) *GetKlinesService {
	return &GetKlinesService{c: c, path: "/fapi/v3/indexPriceKlines", params: map[string]string{
		"pair":     pair,
		"interval": string(interval),
	}}
}

// NewGetMarkPriceKlinesService -- GET /fapi/v3/markPriceKlines
func (c *FuturesClient) NewGetMarkPriceKlinesService(symbol string, interval KlineInterval) *GetKlinesService {
	return &GetKlinesService{c: c, path: "/fapi/v3/markPriceKlines", params: map[string]string{
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

// Do returns parsed klines. Index/Mark price klines return rows with extra
// "ignore" trailing fields; we only populate the meaningful columns.
func (s *GetKlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, s.path, s.params)
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

// GetPremiumIndexService -- GET /fapi/v3/premiumIndex
//
// When symbol is empty the API returns an array; the typed accessor returns
// []PremiumIndex regardless to keep the call site uniform.
type GetPremiumIndexService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetPremiumIndexService() *GetPremiumIndexService {
	return &GetPremiumIndexService{c: c, params: map[string]string{}}
}

func (s *GetPremiumIndexService) SetSymbol(symbol string) *GetPremiumIndexService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetPremiumIndexService) Do(ctx context.Context) ([]PremiumIndex, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/premiumIndex", s.params)
	if _, ok := s.params["symbol"]; ok {
		single, err := request.Do[PremiumIndex](req)
		if err != nil {
			return nil, err
		}
		return []PremiumIndex{*single}, nil
	}
	multi, err := request.Do[[]PremiumIndex](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type PremiumIndex struct {
	Symbol               string          `json:"symbol"`
	MarkPrice            decimal.Decimal `json:"markPrice"`
	IndexPrice           decimal.Decimal `json:"indexPrice"`
	EstimatedSettlePrice decimal.Decimal `json:"estimatedSettlePrice"`
	LastFundingRate      decimal.Decimal `json:"lastFundingRate"`
	NextFundingTime      time.Time       `json:"nextFundingTime,format:unixmilli"`
	InterestRate         decimal.Decimal `json:"interestRate"`
	Time                 time.Time       `json:"time,format:unixmilli"`
}

// GetFundingRateService -- GET /fapi/v3/fundingRate
type GetFundingRateService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetFundingRateService() *GetFundingRateService {
	return &GetFundingRateService{c: c, params: map[string]string{}}
}

func (s *GetFundingRateService) SetSymbol(symbol string) *GetFundingRateService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetFundingRateService) SetStartTime(t time.Time) *GetFundingRateService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFundingRateService) SetEndTime(t time.Time) *GetFundingRateService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFundingRateService) SetLimit(limit int) *GetFundingRateService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFundingRateService) Do(ctx context.Context) ([]FundingRate, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/fundingRate", s.params)
	resp, err := request.Do[[]FundingRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FundingRate struct {
	Symbol      string          `json:"symbol"`
	FundingRate decimal.Decimal `json:"fundingRate"`
	FundingTime time.Time       `json:"fundingTime,format:unixmilli"`
}

// GetFundingInfoService -- GET /fapi/v3/fundingInfo
type GetFundingInfoService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetFundingInfoService() *GetFundingInfoService {
	return &GetFundingInfoService{c: c, params: map[string]string{}}
}

func (s *GetFundingInfoService) SetSymbol(symbol string) *GetFundingInfoService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetFundingInfoService) Do(ctx context.Context) ([]FundingInfo, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/fundingInfo", s.params)
	resp, err := request.Do[[]FundingInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FundingInfo struct {
	Symbol               string          `json:"symbol"`
	InterestRate         decimal.Decimal `json:"interestRate"`
	Time                 time.Time       `json:"time,format:unixmilli"`
	FundingIntervalHours int             `json:"fundingIntervalHours"`
	FundingFeeCap        decimal.Decimal `json:"fundingFeeCap"`
	FundingFeeFloor      decimal.Decimal `json:"fundingFeeFloor"`
}

// Get24hTickerService -- GET /fapi/v3/ticker/24hr
type Get24hTickerService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGet24hTickerService() *Get24hTickerService {
	return &Get24hTickerService{c: c, params: map[string]string{}}
}

func (s *Get24hTickerService) SetSymbol(symbol string) *Get24hTickerService {
	s.params["symbol"] = symbol
	return s
}

func (s *Get24hTickerService) Do(ctx context.Context) ([]Ticker24h, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/ticker/24hr", s.params)
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
}

// GetTickerPriceService -- GET /fapi/v3/ticker/price
type GetTickerPriceService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetTickerPriceService() *GetTickerPriceService {
	return &GetTickerPriceService{c: c, params: map[string]string{}}
}

func (s *GetTickerPriceService) SetSymbol(symbol string) *GetTickerPriceService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetTickerPriceService) Do(ctx context.Context) ([]TickerPrice, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/ticker/price", s.params)
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

// GetBookTickerService -- GET /fapi/v3/ticker/bookTicker
type GetBookTickerService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBookTickerService() *GetBookTickerService {
	return &GetBookTickerService{c: c, params: map[string]string{}}
}

func (s *GetBookTickerService) SetSymbol(symbol string) *GetBookTickerService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetBookTickerService) Do(ctx context.Context) ([]BookTicker, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/ticker/bookTicker", s.params)
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

// GetIndexReferencesService -- GET /fapi/v3/indexreferences
//
// Returns the source-exchange weights used to compute the index price for a
// symbol. Useful for understanding how mark price diverges from spot.
type GetIndexReferencesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetIndexReferencesService(symbol string) *GetIndexReferencesService {
	return &GetIndexReferencesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIndexReferencesService) Do(ctx context.Context) (*IndexReferences, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/indexreferences", s.params)
	return request.Do[IndexReferences](req)
}

type IndexReferences struct {
	Symbol     string           `json:"symbol"`
	Time       time.Time        `json:"time,format:unixmilli"`
	References []IndexReference `json:"references"`
}

type IndexReference struct {
	Exchange string          `json:"exchange"`
	Symbol   string          `json:"symbol"`
	Weight   decimal.Decimal `json:"weight"`
}

// GetOpenInterestService -- GET /fapi/v3/openInterest
//
// Returns the present open interest for a symbol, denominated in the base
// asset (e.g. BTC for BTCUSDT). Multiply by markPrice to get notional in
// the quote asset.
type GetOpenInterestService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOpenInterestService(symbol string) *GetOpenInterestService {
	return &GetOpenInterestService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetOpenInterestService) Do(ctx context.Context) (*OpenInterest, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/openInterest", s.params)
	return request.Do[OpenInterest](req)
}

type OpenInterest struct {
	Symbol       string          `json:"symbol"`
	OpenInterest decimal.Decimal `json:"openInterest"`
	Time         time.Time       `json:"time,format:unixmilli"`
}
