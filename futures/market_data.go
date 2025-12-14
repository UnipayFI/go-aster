package futures

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type DepthService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewDepthService(symbol string) *DepthService {
	return &DepthService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *DepthService) SetLimit(limit int) *DepthService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *DepthService) Do(ctx context.Context) (*DepthResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/depth", s.params)
	return request.Do[DepthResponse](req)
}

type DepthResponse struct {
	LastUpdateId    int64                `json:"lastUpdateId"`
	EventTime       time.Time            `json:"E,format:unixmilli"`
	TransactionTime time.Time            `json:"T,format:unixmilli"`
	Bids            [][2]decimal.Decimal `json:"bids"`
	Asks            [][2]decimal.Decimal `json:"asks"`
}

type GetTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetTradesService(symbol string) *GetTradesService {
	return &GetTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetTradesService) SetLimit(limit int) *GetTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetTradesService) Do(ctx context.Context) ([]TradeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/trades", s.params)
	trades, err := request.Do[[]TradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type TradeResponse struct {
	Id           int64           `json:"id"`
	Price        decimal.Decimal `json:"price"`
	Qty          decimal.Decimal `json:"qty"`
	QuoteQty     decimal.Decimal `json:"quoteQty"`
	Time         time.Time       `json:"time,format:unixmilli"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
}

type HistoricalTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewHistoricalTradesService(symbol string) *HistoricalTradesService {
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
	req := request.Get(ctx, s.c, "/fapi/v1/historicalTrades", s.params).SetApiKeyHeader()
	trades, err := request.Do[[]TradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type AggTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewAggTradesService(symbol string) *AggTradesService {
	return &AggTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *AggTradesService) SetFromId(fromId int64) *AggTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *AggTradesService) SetStartTime(startTime time.Time) *AggTradesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *AggTradesService) SetEndTime(endTime time.Time) *AggTradesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *AggTradesService) SetLimit(limit int) *AggTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *AggTradesService) Do(ctx context.Context) ([]AggTradeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/aggTrades", s.params)
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
	Timestamp    time.Time       `json:"T,format:unixmilli"`
	IsBuyerMaker bool            `json:"m"`
}

type KlinesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewKlinesService(symbol string, interval KlineInterval) *KlinesService {
	return &KlinesService{c: c, params: map[string]string{"symbol": symbol, "interval": string(interval)}}
}

func (s *KlinesService) SetStartTime(startTime time.Time) *KlinesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *KlinesService) SetEndTime(endTime time.Time) *KlinesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *KlinesService) SetLimit(limit int) *KlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *KlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/klines", s.params)
	klines, err := request.Do[[]Kline](req)
	if err != nil {
		return nil, err
	}
	return *klines, nil
}

type Kline struct {
	OpenTime                 time.Time       `json:"openTime,format:unixmilli"`
	Open                     decimal.Decimal `json:"open"`
	High                     decimal.Decimal `json:"high"`
	Low                      decimal.Decimal `json:"low"`
	Close                    decimal.Decimal `json:"close"`
	Volume                   decimal.Decimal `json:"volume"`
	CloseTime                time.Time       `json:"closeTime,format:unixmilli"`
	QuoteAssetVolume         decimal.Decimal `json:"quoteAssetVolume"`
	TradeNum                 int64           `json:"tradeNum"`
	TakerBuyBaseAssetVolume  decimal.Decimal `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume decimal.Decimal `json:"takerBuyQuoteAssetVolume"`
}

func (k *Kline) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal as array: %w", err)
	}
	if len(arr) < 11 {
		return fmt.Errorf("expected array length >= 11, got %d", len(arr))
	}
	var opentime int64
	if err := json.Unmarshal(arr[0], &opentime); err == nil {
		k.OpenTime = time.Unix(opentime/1000, (opentime%1000)*int64(time.Millisecond))
	} else {
		return fmt.Errorf("failed to unmarshal open time: %w", err)
	}
	if err := json.Unmarshal(arr[1], &k.Open); err != nil {
		return fmt.Errorf("failed to unmarshal open: %w", err)
	}
	if err := json.Unmarshal(arr[2], &k.High); err != nil {
		return fmt.Errorf("failed to unmarshal high: %w", err)
	}
	if err := json.Unmarshal(arr[3], &k.Low); err != nil {
		return fmt.Errorf("failed to unmarshal low: %w", err)
	}
	if err := json.Unmarshal(arr[4], &k.Close); err != nil {
		return fmt.Errorf("failed to unmarshal close: %w", err)
	}
	if err := json.Unmarshal(arr[5], &k.Volume); err != nil {
		return fmt.Errorf("failed to unmarshal volume: %w", err)
	}
	var closeTime int64
	if err := json.Unmarshal(arr[6], &closeTime); err == nil {
		k.CloseTime = time.Unix(closeTime/1000, (closeTime%1000)*int64(time.Millisecond))
	} else {
		return fmt.Errorf("failed to unmarshal close time: %w", err)
	}
	if err := json.Unmarshal(arr[7], &k.QuoteAssetVolume); err != nil {
		return fmt.Errorf("failed to unmarshal quote asset volume: %w", err)
	}
	if err := json.Unmarshal(arr[8], &k.TradeNum); err != nil {
		return fmt.Errorf("failed to unmarshal trade num: %w", err)
	}
	if err := json.Unmarshal(arr[9], &k.TakerBuyBaseAssetVolume); err != nil {
		return fmt.Errorf("failed to unmarshal taker buy base asset volume: %w", err)
	}
	if err := json.Unmarshal(arr[10], &k.TakerBuyQuoteAssetVolume); err != nil {
		return fmt.Errorf("failed to unmarshal taker buy quote asset volume: %w", err)
	}
	return nil
}

type IndexPriceKlinesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewIndexPriceKlinesService(pair string, interval KlineInterval) *IndexPriceKlinesService {
	return &IndexPriceKlinesService{c: c, params: map[string]string{"pair": pair, "interval": string(interval)}}
}

func (s *IndexPriceKlinesService) SetStartTime(startTime time.Time) *IndexPriceKlinesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *IndexPriceKlinesService) SetEndTime(endTime time.Time) *IndexPriceKlinesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *IndexPriceKlinesService) SetLimit(limit int) *IndexPriceKlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *IndexPriceKlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/indexPriceKlines", s.params)
	klines, err := request.Do[[]Kline](req)
	if err != nil {
		return nil, err
	}
	return *klines, nil
}

type MarkPriceKlinesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewMarkPriceKlinesService(symbol string, interval KlineInterval) *MarkPriceKlinesService {
	return &MarkPriceKlinesService{c: c, params: map[string]string{"symbol": symbol, "interval": string(interval)}}
}

func (s *MarkPriceKlinesService) SetStartTime(startTime time.Time) *MarkPriceKlinesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *MarkPriceKlinesService) SetEndTime(endTime time.Time) *MarkPriceKlinesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *MarkPriceKlinesService) SetLimit(limit int) *MarkPriceKlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *MarkPriceKlinesService) Do(ctx context.Context) ([]Kline, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/markPriceKlines", s.params)
	klines, err := request.Do[[]Kline](req)
	if err != nil {
		return nil, err
	}
	return *klines, nil
}

type PremiumIndexService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewPremiumIndexService() *PremiumIndexService {
	return &PremiumIndexService{c: c, params: map[string]string{}}
}

func (s *PremiumIndexService) SetSymbol(symbol string) *PremiumIndexService {
	s.params["symbol"] = symbol
	return s
}

func (s *PremiumIndexService) Do(ctx context.Context, symbol string) (*PremiumIndexResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/premiumIndex", map[string]string{"symbol": symbol})
	return request.Do[PremiumIndexResponse](req)
}

func (s *PremiumIndexService) DoAll(ctx context.Context) ([]PremiumIndexResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/premiumIndex", s.params)
	resp, err := request.Do[[]PremiumIndexResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PremiumIndexResponse struct {
	Symbol               string          `json:"symbol"`
	MarkPrice            decimal.Decimal `json:"markPrice"`
	IndexPrice           decimal.Decimal `json:"indexPrice"`
	EstimatedSettlePrice decimal.Decimal `json:"estimatedSettlePrice"`
	LastFundingRate      decimal.Decimal `json:"lastFundingRate"`
	NextFundingTime      time.Time       `json:"nextFundingTime,format:unixmilli"`
	InterestRate         decimal.Decimal `json:"interestRate"`
	Time                 time.Time       `json:"time,format:unixmilli"`
}

type FundingRateService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewFundingRateService() *FundingRateService {
	return &FundingRateService{c: c, params: map[string]string{}}
}

func (s *FundingRateService) SetSymbol(symbol string) *FundingRateService {
	s.params["symbol"] = symbol
	return s
}

func (s *FundingRateService) SetStartTime(startTime time.Time) *FundingRateService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *FundingRateService) SetEndTime(endTime time.Time) *FundingRateService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *FundingRateService) SetLimit(limit int) *FundingRateService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *FundingRateService) Do(ctx context.Context) ([]FundingRateResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/fundingRate", s.params)
	resp, err := request.Do[[]FundingRateResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FundingRateResponse struct {
	Symbol      string          `json:"symbol"`
	FundingRate decimal.Decimal `json:"fundingRate"`
	FundingTime time.Time       `json:"fundingTime,format:unixmilli"`
}

type FundingInfoService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewFundingInfoService() *FundingInfoService {
	return &FundingInfoService{c: c, params: map[string]string{}}
}

func (s *FundingInfoService) SetSymbol(symbol string) *FundingInfoService {
	s.params["symbol"] = symbol
	return s
}

func (s *FundingInfoService) Do(ctx context.Context) ([]FundingInfoResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/fundingInfo", s.params)
	resp, err := request.Do[[]FundingInfoResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FundingInfoResponse struct {
	Symbol               string          `json:"symbol"`
	InterestRate         decimal.Decimal `json:"interestRate"`
	Time                 time.Time       `json:"time,format:unixmilli"`
	FundingIntervalHours int             `json:"fundingIntervalHours"`
	FundingFeeCap        decimal.Decimal `json:"fundingFeeCap"`
	FundingFeeFloor      decimal.Decimal `json:"fundingFeeFloor"`
}

type Ticker24hrService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewTicker24hrService() *Ticker24hrService {
	return &Ticker24hrService{c: c}
}

func (s *Ticker24hrService) Do(ctx context.Context, symbol string) (*Ticker24hrResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/24hr", map[string]string{"symbol": symbol})
	return request.Do[Ticker24hrResponse](req)
}

func (s *Ticker24hrService) DoAll(ctx context.Context) ([]Ticker24hrResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/24hr")
	resp, err := request.Do[[]Ticker24hrResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type Ticker24hrResponse struct {
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

type TickerPriceService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewTickerPriceService() *TickerPriceService {
	return &TickerPriceService{c: c, params: map[string]string{}}
}

func (s *TickerPriceService) SetSymbol(symbol string) *TickerPriceService {
	s.params["symbol"] = symbol
	return s
}

func (s *TickerPriceService) Do(ctx context.Context) (*TickerPriceResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/price", s.params)
	return request.Do[TickerPriceResponse](req)
}

func (s *TickerPriceService) DoAll(ctx context.Context) ([]TickerPriceResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/price", s.params)
	resp, err := request.Do[[]TickerPriceResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type TickerPriceResponse struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Time   time.Time       `json:"time,format:unixmilli"`
}

type BookTickerService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewBookTickerService() *BookTickerService {
	return &BookTickerService{c: c}
}

func (s *BookTickerService) Do(ctx context.Context, symbol string) (*BookTickerResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/bookTicker", map[string]string{"symbol": symbol})
	return request.Do[BookTickerResponse](req)
}

func (s *BookTickerService) DoAll(ctx context.Context) ([]BookTickerResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/ticker/bookTicker")
	resp, err := request.Do[[]BookTickerResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type BookTickerResponse struct {
	Symbol   string          `json:"symbol"`
	BidPrice decimal.Decimal `json:"bidPrice"`
	BidQty   decimal.Decimal `json:"bidQty"`
	AskPrice decimal.Decimal `json:"askPrice"`
	AskQty   decimal.Decimal `json:"askQty"`
	Time     time.Time       `json:"time,format:unixmilli"`
}
