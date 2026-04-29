package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// GetBalanceService -- GET /fapi/v3/balance (USER_DATA)
type GetBalanceService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

func (s *GetBalanceService) Do(ctx context.Context) ([]Balance, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/balance").WithSignature()
	resp, err := request.Do[[]Balance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type Balance struct {
	AccountAlias       string          `json:"accountAlias"`
	Asset              string          `json:"asset"`
	Balance            decimal.Decimal `json:"balance"`
	CrossWalletBalance decimal.Decimal `json:"crossWalletBalance"`
	CrossUnPnl         decimal.Decimal `json:"crossUnPnl"`
	AvailableBalance   decimal.Decimal `json:"availableBalance"`
	MaxWithdrawAmount  decimal.Decimal `json:"maxWithdrawAmount"`
	MarginAvailable    bool            `json:"marginAvailable"`
	UpdateTime         time.Time       `json:"updateTime,format:unixmilli"`
}

// GetAccountService -- GET /fapi/v3/accountWithJoinMargin (USER_DATA)
//
// Note the unusual path: V3 renamed /account to /accountWithJoinMargin to
// signal the response includes joined-margin/multi-asset fields.
type GetAccountService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

func (s *GetAccountService) Do(ctx context.Context) (*AccountInfo, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/accountWithJoinMargin").WithSignature()
	return request.Do[AccountInfo](req)
}

type AccountInfo struct {
	FeeTier                     int             `json:"feeTier"`
	CanTrade                    bool            `json:"canTrade"`
	CanDeposit                  bool            `json:"canDeposit"`
	CanWithdraw                 bool            `json:"canWithdraw"`
	UpdateTime                  int64           `json:"updateTime"`
	TotalInitialMargin          decimal.Decimal `json:"totalInitialMargin"`
	TotalMaintMargin            decimal.Decimal `json:"totalMaintMargin"`
	TotalWalletBalance          decimal.Decimal `json:"totalWalletBalance"`
	TotalUnrealizedProfit       decimal.Decimal `json:"totalUnrealizedProfit"`
	TotalMarginBalance          decimal.Decimal `json:"totalMarginBalance"`
	TotalPositionInitialMargin  decimal.Decimal `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin decimal.Decimal `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     decimal.Decimal `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             decimal.Decimal `json:"totalCrossUnPnl"`
	AvailableBalance            decimal.Decimal `json:"availableBalance"`
	MaxWithdrawAmount           decimal.Decimal `json:"maxWithdrawAmount"`
	Assets                      []AccountAsset  `json:"assets"`
	Positions                   []AccountPosition `json:"positions"`
}

type AccountAsset struct {
	Asset                  string          `json:"asset"`
	WalletBalance          decimal.Decimal `json:"walletBalance"`
	UnrealizedProfit       decimal.Decimal `json:"unrealizedProfit"`
	MarginBalance          decimal.Decimal `json:"marginBalance"`
	MaintMargin            decimal.Decimal `json:"maintMargin"`
	InitialMargin          decimal.Decimal `json:"initialMargin"`
	PositionInitialMargin  decimal.Decimal `json:"positionInitialMargin"`
	OpenOrderInitialMargin decimal.Decimal `json:"openOrderInitialMargin"`
	CrossWalletBalance     decimal.Decimal `json:"crossWalletBalance"`
	CrossUnPnl             decimal.Decimal `json:"crossUnPnl"`
	AvailableBalance       decimal.Decimal `json:"availableBalance"`
	MaxWithdrawAmount      decimal.Decimal `json:"maxWithdrawAmount"`
	MarginAvailable        bool            `json:"marginAvailable"`
	UpdateTime             time.Time       `json:"updateTime,format:unixmilli"`
}

type AccountPosition struct {
	Symbol                 string          `json:"symbol"`
	InitialMargin          decimal.Decimal `json:"initialMargin"`
	MaintMargin            decimal.Decimal `json:"maintMargin"`
	UnrealizedProfit       decimal.Decimal `json:"unrealizedProfit"`
	PositionInitialMargin  decimal.Decimal `json:"positionInitialMargin"`
	OpenOrderInitialMargin decimal.Decimal `json:"openOrderInitialMargin"`
	Leverage               decimal.Decimal `json:"leverage"`
	Isolated               bool            `json:"isolated"`
	EntryPrice             decimal.Decimal `json:"entryPrice"`
	MaxNotional            decimal.Decimal `json:"maxNotional"`
	PositionSide           PositionSide    `json:"positionSide"`
	PositionAmt            decimal.Decimal `json:"positionAmt"`
	UpdateTime             int64           `json:"updateTime"`
}

// GetUserTradesService -- GET /fapi/v3/userTrades (USER_DATA)
type GetUserTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetUserTradesService(symbol string) *GetUserTradesService {
	return &GetUserTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetUserTradesService) SetStartTime(t time.Time) *GetUserTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetEndTime(t time.Time) *GetUserTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetFromId(id int64) *GetUserTradesService {
	s.params["fromId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetUserTradesService) SetLimit(limit int) *GetUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetUserTradesService) Do(ctx context.Context) ([]UserTrade, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/userTrades", s.params).WithSignature()
	resp, err := request.Do[[]UserTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type UserTrade struct {
	Buyer           bool            `json:"buyer"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	ID              int64           `json:"id"`
	Maker           bool            `json:"maker"`
	OrderID         int64           `json:"orderId"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	QuoteQty        decimal.Decimal `json:"quoteQty"`
	RealizedPnl     decimal.Decimal `json:"realizedPnl"`
	Side            OrderSide       `json:"side"`
	PositionSide    PositionSide    `json:"positionSide"`
	Symbol          string          `json:"symbol"`
	Time            time.Time       `json:"time,format:unixmilli"`
}

// GetIncomeHistoryService -- GET /fapi/v3/income (USER_DATA)
type GetIncomeHistoryService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetIncomeHistoryService() *GetIncomeHistoryService {
	return &GetIncomeHistoryService{c: c, params: map[string]string{}}
}

func (s *GetIncomeHistoryService) SetSymbol(symbol string) *GetIncomeHistoryService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetIncomeHistoryService) SetIncomeType(t IncomeType) *GetIncomeHistoryService {
	s.params["incomeType"] = string(t)
	return s
}

func (s *GetIncomeHistoryService) SetStartTime(t time.Time) *GetIncomeHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIncomeHistoryService) SetEndTime(t time.Time) *GetIncomeHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIncomeHistoryService) SetLimit(limit int) *GetIncomeHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIncomeHistoryService) Do(ctx context.Context) ([]IncomeRecord, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/income", s.params).WithSignature()
	resp, err := request.Do[[]IncomeRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type IncomeRecord struct {
	Symbol     string          `json:"symbol"`
	IncomeType IncomeType      `json:"incomeType"`
	Income     decimal.Decimal `json:"income"`
	Asset      string          `json:"asset"`
	Info       string          `json:"info"`
	Time       time.Time       `json:"time,format:unixmilli"`
	TranID     int64           `json:"tranId"`
	TradeID    string          `json:"tradeId"`
}

// GetLeverageBracketService -- GET /fapi/v3/leverageBracket
//
// Without symbol returns []LeverageBracket; with symbol returns a single
// LeverageBracket. Wrap in a slice for uniform handling.
type GetLeverageBracketService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetLeverageBracketService() *GetLeverageBracketService {
	return &GetLeverageBracketService{c: c, params: map[string]string{}}
}

func (s *GetLeverageBracketService) SetSymbol(symbol string) *GetLeverageBracketService {
	s.params["symbol"] = symbol
	return s
}

// Do returns leverage brackets. The server returns a JSON array regardless of
// whether `symbol` is set (with symbol the array has length 1), so we always
// decode as `[]LeverageBracket`.
func (s *GetLeverageBracketService) Do(ctx context.Context) ([]LeverageBracket, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/leverageBracket", s.params).WithSignature()
	multi, err := request.Do[[]LeverageBracket](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type LeverageBracket struct {
	Symbol   string    `json:"symbol"`
	Brackets []Bracket `json:"brackets"`
}

type Bracket struct {
	Bracket          int             `json:"bracket"`
	InitialLeverage  int             `json:"initialLeverage"`
	NotionalCap      decimal.Decimal `json:"notionalCap"`
	NotionalFloor    decimal.Decimal `json:"notionalFloor"`
	MaintMarginRatio decimal.Decimal `json:"maintMarginRatio"`
	Cum              decimal.Decimal `json:"cum"`
}

// GetCommissionRateService -- GET /fapi/v3/commissionRate (USER_DATA)
type GetCommissionRateService struct {
	c      *FuturesClient
	symbol string
}

func (c *FuturesClient) NewGetCommissionRateService(symbol string) *GetCommissionRateService {
	return &GetCommissionRateService{c: c, symbol: symbol}
}

func (s *GetCommissionRateService) Do(ctx context.Context) (*CommissionRateResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/commissionRate", map[string]string{"symbol": s.symbol}).WithSignature()
	return request.Do[CommissionRateResponse](req)
}

type CommissionRateResponse struct {
	Symbol              string          `json:"symbol"`
	MakerCommissionRate decimal.Decimal `json:"makerCommissionRate"`
	TakerCommissionRate decimal.Decimal `json:"takerCommissionRate"`
}

// UpdateMMPService -- POST /fapi/v3/mmp (USER_DATA)
//
// Market Maker Protection: when filled qty / notional / delta in a sliding
// window exceeds the configured limit, MMP-tagged orders are blocked for
// frozenTimeInMilliseconds. windowTime / frozenTime are mandatory; the
// limit triplet (qty/value/delta) is each individually optional.
type UpdateMMPService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewUpdateMMPService(symbol string, windowMs, frozenMs int64) *UpdateMMPService {
	return &UpdateMMPService{c: c, params: map[string]string{
		"symbol":                   symbol,
		"windowTimeInMilliseconds": strconv.FormatInt(windowMs, 10),
		"frozenTimeInMilliseconds": strconv.FormatInt(frozenMs, 10),
	}}
}

func (s *UpdateMMPService) SetQtyLimit(q int64) *UpdateMMPService {
	s.params["qtyLimit"] = strconv.FormatInt(q, 10)
	return s
}

func (s *UpdateMMPService) SetValueLimit(v int64) *UpdateMMPService {
	s.params["valueLimit"] = strconv.FormatInt(v, 10)
	return s
}

func (s *UpdateMMPService) SetDeltaLimit(d int64) *UpdateMMPService {
	s.params["deltaLimit"] = strconv.FormatInt(d, 10)
	return s
}

func (s *UpdateMMPService) Do(ctx context.Context) (bool, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/mmp", s.params).WithSignature()
	resp, err := request.Do[bool](req)
	if err != nil {
		return false, err
	}
	return *resp, nil
}

// GetMMPService -- GET /fapi/v3/mmp (USER_DATA)
type GetMMPService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetMMPService() *GetMMPService {
	return &GetMMPService{c: c, params: map[string]string{}}
}

func (s *GetMMPService) SetSymbol(symbol string) *GetMMPService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMMPService) Do(ctx context.Context) ([]MMPConfig, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/mmp", s.params).WithSignature()
	resp, err := request.Do[[]MMPConfig](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type MMPConfig struct {
	Symbol                   string `json:"symbol"`
	WindowTimeInMilliseconds int64  `json:"windowTimeInMilliseconds"`
	FrozenTimeInMilliseconds int64  `json:"frozenTimeInMilliseconds"`
	QtyLimit                 int64  `json:"qtyLimit"`
	ValueLimit               int64  `json:"valueLimit"`
	DeltaLimit               int64  `json:"deltaLimit"`
}

// DeleteMMPService -- DELETE /fapi/v3/mmp (USER_DATA)
type DeleteMMPService struct {
	c      *FuturesClient
	symbol string
}

func (c *FuturesClient) NewDeleteMMPService(symbol string) *DeleteMMPService {
	return &DeleteMMPService{c: c, symbol: symbol}
}

func (s *DeleteMMPService) Do(ctx context.Context) (bool, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/mmp", map[string]string{"symbol": s.symbol}).WithSignature()
	resp, err := request.Do[bool](req)
	if err != nil {
		return false, err
	}
	return *resp, nil
}

// ResetMMPService -- POST /fapi/v3/mmpReset (USER_DATA)
//
// Manually unfreeze MMP-tagged orders for a symbol after the auto-block
// triggered (instead of waiting for frozenTimeInMilliseconds to elapse).
type ResetMMPService struct {
	c      *FuturesClient
	symbol string
}

func (c *FuturesClient) NewResetMMPService(symbol string) *ResetMMPService {
	return &ResetMMPService{c: c, symbol: symbol}
}

func (s *ResetMMPService) Do(ctx context.Context) (bool, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/mmpReset", map[string]string{"symbol": s.symbol}).WithSignature()
	resp, err := request.Do[bool](req)
	if err != nil {
		return false, err
	}
	return *resp, nil
}
