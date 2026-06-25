package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// ChangePositionModeService -- POST /fapi/v3/positionSide/dual (TRADE)
//
// Switch between One-way Mode (false) and Hedge Mode (true). Applies to
// every symbol on the account.
type ChangePositionModeService struct {
	c                *FuturesClient
	dualSidePosition bool
}

func (c *FuturesClient) NewChangePositionModeService(dual bool) *ChangePositionModeService {
	return &ChangePositionModeService{c: c, dualSidePosition: dual}
}

func (s *ChangePositionModeService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/positionSide/dual", map[string]string{
		"dualSidePosition": strconv.FormatBool(s.dualSidePosition),
	}).WithSignature()
	return request.Do[GenericCodeMsg](req)
}

// GetPositionModeService -- GET /fapi/v3/positionSide/dual (USER_DATA)
type GetPositionModeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetPositionModeService() *GetPositionModeService {
	return &GetPositionModeService{c: c}
}

func (s *GetPositionModeService) Do(ctx context.Context) (*PositionModeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/positionSide/dual").WithSignature()
	return request.Do[PositionModeResponse](req)
}

type PositionModeResponse struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

// ChangeMultiAssetsModeService -- POST /fapi/v3/multiAssetsMargin (TRADE)
type ChangeMultiAssetsModeService struct {
	c                 *FuturesClient
	multiAssetsMargin bool
}

func (c *FuturesClient) NewChangeMultiAssetsModeService(multi bool) *ChangeMultiAssetsModeService {
	return &ChangeMultiAssetsModeService{c: c, multiAssetsMargin: multi}
}

func (s *ChangeMultiAssetsModeService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/multiAssetsMargin", map[string]string{
		"multiAssetsMargin": strconv.FormatBool(s.multiAssetsMargin),
	}).WithSignature()
	return request.Do[GenericCodeMsg](req)
}

// GetMultiAssetsModeService -- GET /fapi/v3/multiAssetsMargin (USER_DATA)
type GetMultiAssetsModeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetMultiAssetsModeService() *GetMultiAssetsModeService {
	return &GetMultiAssetsModeService{c: c}
}

func (s *GetMultiAssetsModeService) Do(ctx context.Context) (*MultiAssetsModeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/multiAssetsMargin").WithSignature()
	return request.Do[MultiAssetsModeResponse](req)
}

type MultiAssetsModeResponse struct {
	MultiAssetsMargin bool `json:"multiAssetsMargin"`
}

// ChangeSTPModeService -- POST /fapi/v3/stpMode (TRADE)
//
// Sets the account-level Self-Trade Prevention mode applied to every symbol by
// default; individual orders may override it via PlaceOrderService.SetSTPMode.
type ChangeSTPModeService struct {
	c       *FuturesClient
	stpMode STPMode
}

func (c *FuturesClient) NewChangeSTPModeService(mode STPMode) *ChangeSTPModeService {
	return &ChangeSTPModeService{c: c, stpMode: mode}
}

func (s *ChangeSTPModeService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/stpMode", map[string]string{
		"stpMode": string(s.stpMode),
	}).WithSignature()
	return request.Do[GenericCodeMsg](req)
}

// GetSTPModeService -- GET /fapi/v3/stpMode (USER_DATA)
type GetSTPModeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetSTPModeService() *GetSTPModeService {
	return &GetSTPModeService{c: c}
}

func (s *GetSTPModeService) Do(ctx context.Context) (*STPModeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/stpMode").WithSignature()
	return request.Do[STPModeResponse](req)
}

type STPModeResponse struct {
	STPMode STPMode `json:"stpMode"`
}

// ChangeLeverageService -- POST /fapi/v3/leverage (TRADE)
type ChangeLeverageService struct {
	c        *FuturesClient
	symbol   string
	leverage int
}

func (c *FuturesClient) NewChangeLeverageService(symbol string, leverage int) *ChangeLeverageService {
	return &ChangeLeverageService{c: c, symbol: symbol, leverage: leverage}
}

func (s *ChangeLeverageService) Do(ctx context.Context) (*ChangeLeverageResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/leverage", map[string]string{
		"symbol":   s.symbol,
		"leverage": strconv.Itoa(s.leverage),
	}).WithSignature()
	return request.Do[ChangeLeverageResponse](req)
}

type ChangeLeverageResponse struct {
	Leverage         int             `json:"leverage"`
	MaxNotionalValue decimal.Decimal `json:"maxNotionalValue"`
	Symbol           string          `json:"symbol"`
}

// ChangeMarginTypeService -- POST /fapi/v3/marginType (TRADE)
type ChangeMarginTypeService struct {
	c          *FuturesClient
	symbol     string
	marginType MarginType
}

func (c *FuturesClient) NewChangeMarginTypeService(symbol string, m MarginType) *ChangeMarginTypeService {
	return &ChangeMarginTypeService{c: c, symbol: symbol, marginType: m}
}

func (s *ChangeMarginTypeService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/marginType", map[string]string{
		"symbol":     s.symbol,
		"marginType": string(s.marginType),
	}).WithSignature()
	return request.Do[GenericCodeMsg](req)
}

// ModifyIsolatedPositionMarginService -- POST /fapi/v3/positionMargin (TRADE)
type ModifyIsolatedPositionMarginService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewModifyIsolatedPositionMarginService(symbol string, amount decimal.Decimal, action PositionMarginType) *ModifyIsolatedPositionMarginService {
	return &ModifyIsolatedPositionMarginService{c: c, params: map[string]string{
		"symbol": symbol,
		"amount": amount.String(),
		"type":   strconv.Itoa(int(action)),
	}}
}

func (s *ModifyIsolatedPositionMarginService) SetPositionSide(p PositionSide) *ModifyIsolatedPositionMarginService {
	s.params["positionSide"] = string(p)
	return s
}

func (s *ModifyIsolatedPositionMarginService) Do(ctx context.Context) (*ModifyMarginResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/positionMargin", s.params).WithSignature()
	return request.Do[ModifyMarginResponse](req)
}

type ModifyMarginResponse struct {
	Amount decimal.Decimal `json:"amount"`
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Type   int             `json:"type"`
}

// GetPositionMarginHistoryService -- GET /fapi/v3/positionMargin/history
type GetPositionMarginHistoryService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetPositionMarginHistoryService(symbol string) *GetPositionMarginHistoryService {
	return &GetPositionMarginHistoryService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetPositionMarginHistoryService) SetType(t PositionMarginType) *GetPositionMarginHistoryService {
	s.params["type"] = strconv.Itoa(int(t))
	return s
}

func (s *GetPositionMarginHistoryService) SetStartTime(t time.Time) *GetPositionMarginHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetPositionMarginHistoryService) SetEndTime(t time.Time) *GetPositionMarginHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetPositionMarginHistoryService) SetLimit(limit int) *GetPositionMarginHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetPositionMarginHistoryService) Do(ctx context.Context) ([]PositionMarginHistoryEntry, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/positionMargin/history", s.params).WithSignature()
	resp, err := request.Do[[]PositionMarginHistoryEntry](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PositionMarginHistoryEntry struct {
	Amount       decimal.Decimal `json:"amount"`
	Asset        string          `json:"asset"`
	Symbol       string          `json:"symbol"`
	Time         time.Time       `json:"time,format:unixmilli"`
	Type         int             `json:"type"`
	PositionSide PositionSide    `json:"positionSide"`
}

// GetPositionRiskService -- GET /fapi/v3/positionRisk (USER_DATA)
type GetPositionRiskService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetPositionRiskService() *GetPositionRiskService {
	return &GetPositionRiskService{c: c, params: map[string]string{}}
}

func (s *GetPositionRiskService) SetSymbol(symbol string) *GetPositionRiskService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetPositionRiskService) Do(ctx context.Context) ([]PositionRisk, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/positionRisk", s.params).WithSignature()
	resp, err := request.Do[[]PositionRisk](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PositionRisk struct {
	EntryPrice       decimal.Decimal `json:"entryPrice"`
	MarginType       string          `json:"marginType"` // "isolated" / "cross" -- lowercase wire format
	IsAutoAddMargin  string          `json:"isAutoAddMargin"`
	IsolatedMargin   decimal.Decimal `json:"isolatedMargin"`
	Leverage         decimal.Decimal `json:"leverage"`
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"`
	MarkPrice        decimal.Decimal `json:"markPrice"`
	MaxNotionalValue decimal.Decimal `json:"maxNotionalValue"`
	PositionAmt      decimal.Decimal `json:"positionAmt"`
	Symbol           string          `json:"symbol"`
	UnRealizedProfit decimal.Decimal `json:"unRealizedProfit"`
	PositionSide     PositionSide    `json:"positionSide"`
	UpdateTime       time.Time       `json:"updateTime,format:unixmilli"`
}

// GetADLQuantileService -- GET /fapi/v3/adlQuantile (USER_DATA)
type GetADLQuantileService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetADLQuantileService() *GetADLQuantileService {
	return &GetADLQuantileService{c: c, params: map[string]string{}}
}

func (s *GetADLQuantileService) SetSymbol(symbol string) *GetADLQuantileService {
	s.params["symbol"] = symbol
	return s
}

// Do returns ADL quantiles. With symbol the server returns a single object;
// without it, an array. We normalize both shapes to []ADLQuantile.
func (s *GetADLQuantileService) Do(ctx context.Context) ([]ADLQuantile, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/adlQuantile", s.params).WithSignature()
	if _, ok := s.params["symbol"]; ok {
		single, err := request.Do[ADLQuantile](req)
		if err != nil {
			return nil, err
		}
		return []ADLQuantile{*single}, nil
	}
	multi, err := request.Do[[]ADLQuantile](req)
	if err != nil {
		return nil, err
	}
	return *multi, nil
}

type ADLQuantile struct {
	Symbol      string         `json:"symbol"`
	AdlQuantile map[string]int `json:"adlQuantile"` // keys: "LONG"/"SHORT"/"BOTH"/"HEDGE"
}

// GetForceOrdersService -- GET /fapi/v3/forceOrders (USER_DATA)
type GetForceOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetForceOrdersService() *GetForceOrdersService {
	return &GetForceOrdersService{c: c, params: map[string]string{}}
}

func (s *GetForceOrdersService) SetSymbol(symbol string) *GetForceOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetForceOrdersService) SetAutoCloseType(t AutoCloseType) *GetForceOrdersService {
	s.params["autoCloseType"] = string(t)
	return s
}

func (s *GetForceOrdersService) SetStartTime(t time.Time) *GetForceOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetForceOrdersService) SetEndTime(t time.Time) *GetForceOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetForceOrdersService) SetLimit(limit int) *GetForceOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetForceOrdersService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/forceOrders", s.params).WithSignature()
	resp, err := request.Do[[]Order](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
