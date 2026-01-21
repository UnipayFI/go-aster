package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
)

type ChangePositionModeService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewChangePositionModeService(dualSidePosition bool) *ChangePositionModeService {
	return &ChangePositionModeService{
		c:      c,
		params: map[string]string{"dualSidePosition": strconv.FormatBool(dualSidePosition)},
	}
}

func (s *ChangePositionModeService) Do(ctx context.Context) error {
	req := request.Post(s.c, ctx, "/fapi/v1/positionSide/dual", s.params).WithSignature()
	return handlerGeneralResponse(request.Do[GeneralResponse](req))
}

type GetPositionModeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetPositionModeService() *GetPositionModeService {
	return &GetPositionModeService{c: c}
}

func (s *GetPositionModeService) Do(ctx context.Context) (dualSidePosition bool, err error) {
	req := request.Get(ctx, s.c, "/fapi/v1/positionSide/dual").WithSignature()
	resp, err := request.Do[struct {
		DualSidePosition bool `json:"dualSidePosition"`
	}](req)
	if err != nil {
		return false, err
	}
	return resp.DualSidePosition, nil
}

type ChangeMultiAssetsModeService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewChangeMultiAssetsModeService(multiAssetsMargin bool) *ChangeMultiAssetsModeService {
	return &ChangeMultiAssetsModeService{
		c:      c,
		params: map[string]string{"multiAssetsMargin": strconv.FormatBool(multiAssetsMargin)},
	}
}

func (s *ChangeMultiAssetsModeService) Do(ctx context.Context) error {
	req := request.Post(s.c, ctx, "/fapi/v1/multiAssetsMargin", s.params).WithSignature()
	return handlerGeneralResponse(request.Do[GeneralResponse](req))
}

type GetMultiAssetsModeService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetMultiAssetsModeService() *GetMultiAssetsModeService {
	return &GetMultiAssetsModeService{c: c}
}

func (s *GetMultiAssetsModeService) Do(ctx context.Context) (multiAssetsMargin bool, err error) {
	req := request.Get(ctx, s.c, "/fapi/v1/multiAssetsMargin").WithSignature()
	resp, err := request.Do[struct {
		MultiAssetsMargin bool `json:"multiAssetsMargin"`
	}](req)
	if err != nil {
		return false, err
	}
	return resp.MultiAssetsMargin, nil
}

type ChangeLeverageService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewChangeLeverageService(symbol string, leverage int) *ChangeLeverageService {
	return &ChangeLeverageService{
		c: c,
		params: map[string]string{
			"symbol":   symbol,
			"leverage": strconv.Itoa(leverage),
		},
	}
}

func (s *ChangeLeverageService) Do(ctx context.Context) (*ChangeLeverageResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/leverage", s.params).WithSignature()
	return request.Do[ChangeLeverageResponse](req)
}

type ChangeLeverageResponse struct {
	Leverage         int             `json:"leverage"`
	MaxNotionalValue decimal.Decimal `json:"maxNotionalValue"`
	Symbol           string          `json:"symbol"`
}

type ChangeMarginTypeService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewChangeMarginTypeService(symbol string, marginType MarginType) *ChangeMarginTypeService {
	return &ChangeMarginTypeService{
		c: c,
		params: map[string]string{
			"symbol":     symbol,
			"marginType": string(marginType),
		},
	}
}

func (s *ChangeMarginTypeService) Do(ctx context.Context) error {
	req := request.Post(s.c, ctx, "/fapi/v1/marginType", s.params).WithSignature()
	return handlerGeneralResponse(request.Do[GeneralResponse](req))
}

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

func (s *GetPositionRiskService) Do(ctx context.Context) ([]PositionRiskResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v2/positionRisk", s.params).WithSignature()
	resp, err := request.Do[[]PositionRiskResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PositionRiskResponse struct {
	Symbol           string          `json:"symbol"`
	PositionAmt      decimal.Decimal `json:"positionAmt"`
	EntryPrice       decimal.Decimal `json:"entryPrice"`
	MarkPrice        decimal.Decimal `json:"markPrice"`
	UnRealizedProfit decimal.Decimal `json:"unRealizedProfit"`
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"`
	Leverage         int             `json:"leverage,string"`
	MaxNotionalValue decimal.Decimal `json:"maxNotionalValue"`
	MarginType       MarginType      `json:"marginType"`
	IsolatedMargin   decimal.Decimal `json:"isolatedMargin"`
	IsAutoAddMargin  string          `json:"isAutoAddMargin"`
	PositionSide     PositionSide    `json:"positionSide"`
	UpdateTime       time.Time       `json:"updateTime,format:unixmilli"`
}

type AddIsolatedMarginService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewAddIsolatedMarginService(symbol string, amount decimal.Decimal, marginType int) *AddIsolatedMarginService {
	return &AddIsolatedMarginService{
		c: c,
		params: map[string]string{
			"symbol": symbol,
			"amount": amount.String(),
			"type":   strconv.Itoa(marginType),
		},
	}
}

func (s *AddIsolatedMarginService) SetPositionSide(positionSide PositionSide) *AddIsolatedMarginService {
	s.params["positionSide"] = string(positionSide)
	return s
}

func (s *AddIsolatedMarginService) Do(ctx context.Context) (*AddIsolatedMarginResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/positionMargin", s.params).WithSignature()
	return request.Do[AddIsolatedMarginResponse](req)
}

type AddIsolatedMarginResponse struct {
	Amount decimal.Decimal `json:"amount"`
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Type   int             `json:"type"`
}

type GetBalanceService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

func (s *GetBalanceService) Do(ctx context.Context) ([]BalanceResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v2/balance").WithSignature()
	resp, err := request.Do[[]BalanceResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type BalanceResponse struct {
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

type GetAccountService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

func (s *GetAccountService) Do(ctx context.Context) (*AccountResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v4/account").WithSignature()
	return request.Do[AccountResponse](req)
}

type AccountResponse struct {
	FeeTier                     int               `json:"feeTier"`
	CanTrade                    bool              `json:"canTrade"`
	CanDeposit                  bool              `json:"canDeposit"`
	CanWithdraw                 bool              `json:"canWithdraw"`
	UpdateTime                  time.Time         `json:"updateTime,format:unixmilli"`
	TotalInitialMargin          decimal.Decimal   `json:"totalInitialMargin"`
	TotalMaintMargin            decimal.Decimal   `json:"totalMaintMargin"`
	TotalWalletBalance          decimal.Decimal   `json:"totalWalletBalance"`
	TotalUnrealizedProfit       decimal.Decimal   `json:"totalUnrealizedProfit"`
	TotalMarginBalance          decimal.Decimal   `json:"totalMarginBalance"`
	TotalPositionInitialMargin  decimal.Decimal   `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin decimal.Decimal   `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     decimal.Decimal   `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             decimal.Decimal   `json:"totalCrossUnPnl"`
	AvailableBalance            decimal.Decimal   `json:"availableBalance"`
	MaxWithdrawAmount           decimal.Decimal   `json:"maxWithdrawAmount"`
	Assets                      []AccountAsset    `json:"assets"`
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
	Leverage               int             `json:"leverage,string"`
	Isolated               bool            `json:"isolated"`
	EntryPrice             decimal.Decimal `json:"entryPrice"`
	MaxNotional            decimal.Decimal `json:"maxNotional"`
	PositionSide           PositionSide    `json:"positionSide"`
	PositionAmt            decimal.Decimal `json:"positionAmt"`
	UpdateTime             time.Time       `json:"updateTime,format:unixmilli"`
}

type GetIncomeService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetIncomeService() *GetIncomeService {
	return &GetIncomeService{c: c, params: map[string]string{}}
}

func (s *GetIncomeService) SetSymbol(symbol string) *GetIncomeService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetIncomeService) SetIncomeType(incomeType IncomeType) *GetIncomeService {
	s.params["incomeType"] = string(incomeType)
	return s
}

func (s *GetIncomeService) SetStartTime(startTime time.Time) *GetIncomeService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetIncomeService) SetEndTime(endTime time.Time) *GetIncomeService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetIncomeService) SetLimit(limit int) *GetIncomeService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIncomeService) Do(ctx context.Context) ([]IncomeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/income", s.params).WithSignature()
	resp, err := request.Do[[]IncomeResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type IncomeResponse struct {
	Symbol     string          `json:"symbol"`
	IncomeType IncomeType      `json:"incomeType"`
	Income     decimal.Decimal `json:"income"`
	Asset      string          `json:"asset"`
	Info       string          `json:"info"`
	Time       time.Time       `json:"time,format:unixmilli"`
	TranId     int64           `json:"tranId"`
	TradeId    string          `json:"tradeId"`
}

type GetUserTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetUserTradesService(symbol string) *GetUserTradesService {
	return &GetUserTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetUserTradesService) SetStartTime(startTime time.Time) *GetUserTradesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetEndTime(endTime time.Time) *GetUserTradesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetFromId(fromId int64) *GetUserTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *GetUserTradesService) SetLimit(limit int) *GetUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetUserTradesService) Do(ctx context.Context) ([]UserTradeResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/userTrades", s.params).WithSignature()
	resp, err := request.Do[[]UserTradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type UserTradeResponse struct {
	Symbol          string          `json:"symbol"`
	Id              int64           `json:"id"`
	OrderId         int64           `json:"orderId"`
	Side            OrderSide       `json:"side"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	RealizedPnl     decimal.Decimal `json:"realizedPnl"`
	MarginAsset     string          `json:"marginAsset"`
	QuoteQty        decimal.Decimal `json:"quoteQty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	Time            time.Time       `json:"time,format:unixmilli"`
	PositionSide    PositionSide    `json:"positionSide"`
	Maker           bool            `json:"maker"`
	Buyer           bool            `json:"buyer"`
}

type GetCommissionRateService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetCommissionRateService(symbol string) *GetCommissionRateService {
	return &GetCommissionRateService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetCommissionRateService) Do(ctx context.Context) (*CommissionRateResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/commissionRate", s.params).WithSignature()
	return request.Do[CommissionRateResponse](req)
}

type CommissionRateResponse struct {
	Symbol              string          `json:"symbol"`
	MakerCommissionRate decimal.Decimal `json:"makerCommissionRate"`
	TakerCommissionRate decimal.Decimal `json:"takerCommissionRate"`
}

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

func (s *GetForceOrdersService) SetAutoCloseType(autoCloseType AutoCloseType) *GetForceOrdersService {
	s.params["autoCloseType"] = string(autoCloseType)
	return s
}

func (s *GetForceOrdersService) SetStartTime(startTime time.Time) *GetForceOrdersService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetForceOrdersService) SetEndTime(endTime time.Time) *GetForceOrdersService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetForceOrdersService) SetLimit(limit int) *GetForceOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetForceOrdersService) Do(ctx context.Context) ([]ForceOrderResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/forceOrders", s.params).WithSignature()
	resp, err := request.Do[[]ForceOrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type ForceOrderResponse struct {
	OrderId       int64           `json:"orderId"`
	Symbol        string          `json:"symbol"`
	Status        OrderStatus     `json:"status"`
	ClientOrderId string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	OrigQty       decimal.Decimal `json:"origQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	CumQuote      decimal.Decimal `json:"cumQuote"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	Type          OrderType       `json:"type"`
	ReduceOnly    bool            `json:"reduceOnly"`
	ClosePosition bool            `json:"closePosition"`
	Side          OrderSide       `json:"side"`
	PositionSide  PositionSide    `json:"positionSide"`
	StopPrice     decimal.Decimal `json:"stopPrice"`
	WorkingType   WorkingType     `json:"workingType"`
	OrigType      OrderType       `json:"origType"`
	Time          time.Time       `json:"time,format:unixmilli"`
	UpdateTime    time.Time       `json:"updateTime,format:unixmilli"`
}

type GetPositionMarginHistoryService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetPositionMarginHistoryService(symbol string) *GetPositionMarginHistoryService {
	return &GetPositionMarginHistoryService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetPositionMarginHistoryService) SetType(marginType int) *GetPositionMarginHistoryService {
	s.params["type"] = strconv.Itoa(marginType)
	return s
}

func (s *GetPositionMarginHistoryService) SetStartTime(startTime time.Time) *GetPositionMarginHistoryService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetPositionMarginHistoryService) SetEndTime(endTime time.Time) *GetPositionMarginHistoryService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetPositionMarginHistoryService) SetLimit(limit int) *GetPositionMarginHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetPositionMarginHistoryService) Do(ctx context.Context) ([]PositionMarginHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/positionMargin/history", s.params).WithSignature()
	resp, err := request.Do[[]PositionMarginHistoryResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PositionMarginHistoryResponse struct {
	Amount       decimal.Decimal `json:"amount"`
	Asset        string          `json:"asset"`
	Symbol       string          `json:"symbol"`
	Time         time.Time       `json:"time,format:unixmilli"`
	Type         int             `json:"type"`
	PositionSide PositionSide    `json:"positionSide"`
}

type GetLeverageBracketService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetLeverageBracketService() *GetLeverageBracketService {
	return &GetLeverageBracketService{c: c}
}

func (s *GetLeverageBracketService) DoAll(ctx context.Context) ([]LeverageBracketResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/leverageBracket").WithSignature()
	resp, err := request.Do[[]LeverageBracketResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

func (s *GetLeverageBracketService) Do(ctx context.Context, symbol string) ([]LeverageBracketResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/leverageBracket", map[string]string{"symbol": symbol}).WithSignature()
	resp, err := request.Do[[]LeverageBracketResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type LeverageBracketResponse struct {
	Symbol   string            `json:"symbol"`
	Brackets []LeverageBracket `json:"brackets"`
}

type LeverageBracket struct {
	Bracket          int             `json:"bracket"`
	InitialLeverage  int             `json:"initialLeverage"`
	NotionalCap      decimal.Decimal `json:"notionalCap"`
	NotionalFloor    decimal.Decimal `json:"notionalFloor"`
	MaintMarginRatio decimal.Decimal `json:"maintMarginRatio"`
	Cum              decimal.Decimal `json:"cum"`
}

type GetAdlQuantileService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetAdlQuantileService() *GetAdlQuantileService {
	return &GetAdlQuantileService{c: c}
}

func (s *GetAdlQuantileService) DoAll(ctx context.Context) ([]AdlQuantileResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/adlQuantile").WithSignature()
	resp, err := request.Do[[]AdlQuantileResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

func (s *GetAdlQuantileService) Do(ctx context.Context, symbol string) (*AdlQuantileResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/adlQuantile", map[string]string{"symbol": symbol}).WithSignature()
	return request.Do[AdlQuantileResponse](req)
}

type AdlQuantileResponse struct {
	Symbol      string      `json:"symbol"`
	AdlQuantile AdlQuantile `json:"adlQuantile"`
}

type AdlQuantile struct {
	Long  int `json:"LONG"`
	Short int `json:"SHORT"`
	Both  int `json:"BOTH"`
	Hedge int `json:"HEDGE"`
}

type CreateListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewCreateListenKeyService() *CreateListenKeyService {
	return &CreateListenKeyService{c: c}
}

func (s *CreateListenKeyService) Do(ctx context.Context) (string, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/listenKey").WithApiKey()
	resp, err := request.Do[struct {
		ListenKey string `json:"listenKey"`
	}](req)
	if err != nil {
		return "", err
	}
	return resp.ListenKey, nil
}

type ExtendListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewExtendListenKeyService() *ExtendListenKeyService {
	return &ExtendListenKeyService{c: c}
}

func (s *ExtendListenKeyService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/fapi/v1/listenKey").WithApiKey()
	_, err := request.Do[struct{}](req)
	return err
}

type CloseListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewCloseListenKeyService() *CloseListenKeyService {
	return &CloseListenKeyService{c: c}
}

func (s *CloseListenKeyService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/fapi/v1/listenKey").WithApiKey()
	_, err := request.Do[struct{}](req)
	return err
}

type RemainingOpenableNotionalValueService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewRemainingOpenableNotionalValueService(symbol string, leverage int) *RemainingOpenableNotionalValueService {
	return &RemainingOpenableNotionalValueService{c: c, params: map[string]string{"symbol": symbol, "leverage": strconv.Itoa(leverage)}}
}

func (s *RemainingOpenableNotionalValueService) Do(ctx context.Context) (float64, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/remainingOpenableNotionalValue", s.params).WithSignature()
	resp, err := request.Do[struct {
		RemainingOpenableNotionalValue float64 `json:"remainingOpenableNotionalValue"`
	}](req)
	if err != nil {
		return 0, err
	}
	return resp.RemainingOpenableNotionalValue, nil
}
