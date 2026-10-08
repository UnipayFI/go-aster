package futures

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// The builder user queries share their request and response shape: the
// authenticated account is the builder identity, and userAddresses narrows the
// query to at most 50 users who approved the caller as their builder (others
// come back in Errors) instead of every bound user. For the account, open
// order, balance, position and commission rate queries below, page/limit walk
// the bound users in binding order, and the caller must itself have deposited
// at least once.

// GetBuilderUserAccountsService -- GET /fapi/v3/builder/userAccounts (USER_DATA)
//
// Queries account information for the users trading under the caller's
// builder code.
type GetBuilderUserAccountsService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserAccountsService() *GetBuilderUserAccountsService {
	return &GetBuilderUserAccountsService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns these users (at most 50).
func (s *GetBuilderUserAccountsService) SetUserAddresses(addresses ...string) *GetBuilderUserAccountsService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetSymbol only returns this symbol's entry in SymbolConfig.
func (s *GetBuilderUserAccountsService) SetSymbol(symbol string) *GetBuilderUserAccountsService {
	s.params["symbol"] = symbol
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserAccountsService) SetPage(page int) *GetBuilderUserAccountsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of users per page (default 50, maximum 1000).
func (s *GetBuilderUserAccountsService) SetLimit(limit int) *GetBuilderUserAccountsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserAccountsService) Do(ctx context.Context) (*BuilderUserAccounts, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userAccounts", s.params).WithSignature()
	return request.Do[BuilderUserAccounts](req)
}

// BuilderUserAccounts is the paginated response of GetBuilderUserAccountsService.
type BuilderUserAccounts struct {
	Total       int64                `json:"total"`
	CurrentPage int                  `json:"currentPage"`
	TotalPages  int                  `json:"totalPages"`
	PageSize    int                  `json:"pageSize"`
	HasMore     bool                 `json:"hasMore"`
	Rows        []BuilderUserAccount `json:"rows"`
	Errors      []BuilderUserError   `json:"errors"`
}

// BuilderUserAccount is the account of one user under the caller's builder
// code. JointMargin reports Multi-Assets Mode.
type BuilderUserAccount struct {
	Address          string                    `json:"address"`
	FeeTier          int                       `json:"feeTier"`
	CanTrade         bool                      `json:"canTrade"`
	CanDeposit       bool                      `json:"canDeposit"`
	CanWithdraw      bool                      `json:"canWithdraw"`
	UpdateTime       time.Time                 `json:"updateTime,format:unixmilli"`
	AccountType      int                       `json:"accountType"`
	DualSidePosition bool                      `json:"dualSidePosition"`
	JointMargin      bool                      `json:"jointMargin"`
	FeeBurn          bool                      `json:"feeBurn"`
	FeeBurnAssetID   int                       `json:"feeBurnAssetId"`
	SymbolConfig     []BuilderUserSymbolConfig `json:"symbolConfig"`
}

// BuilderUserSymbolConfig is a user's per-symbol configuration.
type BuilderUserSymbolConfig struct {
	Symbol            string          `json:"symbol"`
	Leverage          int             `json:"leverage"`
	NotionalLimitCoef decimal.Decimal `json:"notionalLimitCoef"`
}

// GetBuilderUserOpenOrdersService -- GET /fapi/v3/builder/userOpenOrders (USER_DATA)
//
// Queries the current open orders of the users trading under the caller's
// builder code. page/limit page through users, not orders, so each user may
// contribute any number of rows.
type GetBuilderUserOpenOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserOpenOrdersService() *GetBuilderUserOpenOrdersService {
	return &GetBuilderUserOpenOrdersService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns orders of these users (at most 50).
func (s *GetBuilderUserOpenOrdersService) SetUserAddresses(addresses ...string) *GetBuilderUserOpenOrdersService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetSymbol only returns open orders on this symbol.
func (s *GetBuilderUserOpenOrdersService) SetSymbol(symbol string) *GetBuilderUserOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserOpenOrdersService) SetPage(page int) *GetBuilderUserOpenOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of users per page (default 50, maximum 1000).
func (s *GetBuilderUserOpenOrdersService) SetLimit(limit int) *GetBuilderUserOpenOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserOpenOrdersService) Do(ctx context.Context) (*BuilderUserOpenOrders, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userOpenOrders", s.params).WithSignature()
	return request.Do[BuilderUserOpenOrders](req)
}

// BuilderUserOpenOrders is the paginated response of
// GetBuilderUserOpenOrdersService.
type BuilderUserOpenOrders struct {
	Total       int64                  `json:"total"`
	CurrentPage int                    `json:"currentPage"`
	TotalPages  int                    `json:"totalPages"`
	PageSize    int                    `json:"pageSize"`
	HasMore     bool                   `json:"hasMore"`
	Rows        []BuilderUserOpenOrder `json:"rows"`
	Errors      []BuilderUserError     `json:"errors"`
}

// BuilderUserOpenOrder is one open order of a user under the caller's builder
// code.
type BuilderUserOpenOrder struct {
	Address       string          `json:"address"`
	Symbol        string          `json:"symbol"`
	OrderID       int64           `json:"orderId"`
	ClientOrderID string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	OrigQty       decimal.Decimal `json:"origQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	StopPrice     decimal.Decimal `json:"stopPrice"`
	Status        OrderStatus     `json:"status"`
	Side          OrderSide       `json:"side"`
	PositionSide  PositionSide    `json:"positionSide"`
	Type          OrderType       `json:"type"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	Time          time.Time       `json:"time,format:unixmilli"`
	WorkingType   WorkingType     `json:"workingType"`
}

// GetBuilderUserBalancesService -- GET /fapi/v3/builder/userBalances (USER_DATA)
//
// Queries the wallet balances of the users trading under the caller's builder
// code.
type GetBuilderUserBalancesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserBalancesService() *GetBuilderUserBalancesService {
	return &GetBuilderUserBalancesService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns balances of these users (at most 50).
func (s *GetBuilderUserBalancesService) SetUserAddresses(addresses ...string) *GetBuilderUserBalancesService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserBalancesService) SetPage(page int) *GetBuilderUserBalancesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of users per page (default 50, maximum 1000).
func (s *GetBuilderUserBalancesService) SetLimit(limit int) *GetBuilderUserBalancesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserBalancesService) Do(ctx context.Context) (*BuilderUserBalances, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userBalances", s.params).WithSignature()
	return request.Do[BuilderUserBalances](req)
}

// BuilderUserBalances is the paginated response of
// GetBuilderUserBalancesService.
type BuilderUserBalances struct {
	Total       int64                `json:"total"`
	CurrentPage int                  `json:"currentPage"`
	TotalPages  int                  `json:"totalPages"`
	PageSize    int                  `json:"pageSize"`
	HasMore     bool                 `json:"hasMore"`
	Rows        []BuilderUserBalance `json:"rows"`
	Errors      []BuilderUserError   `json:"errors"`
}

// BuilderUserBalance is one asset balance of a user under the caller's builder
// code; BalanceInUSD is WalletBalance at Price.
type BuilderUserBalance struct {
	Address       string          `json:"address"`
	Asset         string          `json:"asset"`
	WalletBalance decimal.Decimal `json:"walletBalance"`
	Price         decimal.Decimal `json:"price"`
	BalanceInUSD  decimal.Decimal `json:"balanceInUsd"`
}

// GetBuilderUserPositionRiskService -- GET /fapi/v3/builder/userPositionRisk (USER_DATA)
//
// Queries the current positions of the users trading under the caller's
// builder code. Unlike GetPositionRiskService, only positions with a non-zero
// positionAmt are returned: the BOTH side for One-way Mode users, LONG/SHORT
// for Hedge Mode users. page/limit page through users, not positions.
type GetBuilderUserPositionRiskService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserPositionRiskService() *GetBuilderUserPositionRiskService {
	return &GetBuilderUserPositionRiskService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns positions of these users (at most 50).
func (s *GetBuilderUserPositionRiskService) SetUserAddresses(addresses ...string) *GetBuilderUserPositionRiskService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetSymbol only returns positions on this symbol.
func (s *GetBuilderUserPositionRiskService) SetSymbol(symbol string) *GetBuilderUserPositionRiskService {
	s.params["symbol"] = symbol
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserPositionRiskService) SetPage(page int) *GetBuilderUserPositionRiskService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of users per page (default 50, maximum 1000).
func (s *GetBuilderUserPositionRiskService) SetLimit(limit int) *GetBuilderUserPositionRiskService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserPositionRiskService) Do(ctx context.Context) (*BuilderUserPositionRisk, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userPositionRisk", s.params).WithSignature()
	return request.Do[BuilderUserPositionRisk](req)
}

// BuilderUserPositionRisk is the paginated response of
// GetBuilderUserPositionRiskService.
type BuilderUserPositionRisk struct {
	Total       int64                 `json:"total"`
	CurrentPage int                   `json:"currentPage"`
	TotalPages  int                   `json:"totalPages"`
	PageSize    int                   `json:"pageSize"`
	HasMore     bool                  `json:"hasMore"`
	Rows        []BuilderUserPosition `json:"rows"`
	Errors      []BuilderUserError    `json:"errors"`
}

// BuilderUserPosition is one position of a user under the caller's builder
// code. LiquidationPrice is 0 whenever the underlying value would be negative.
type BuilderUserPosition struct {
	Address          string          `json:"address"`
	Symbol           string          `json:"symbol"`
	PositionAmt      decimal.Decimal `json:"positionAmt"`
	EntryPrice       decimal.Decimal `json:"entryPrice"`
	MarkPrice        decimal.Decimal `json:"markPrice"`
	UnRealizedProfit decimal.Decimal `json:"unRealizedProfit"`
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"`
	Leverage         decimal.Decimal `json:"leverage"`
	MaxNotionalValue decimal.Decimal `json:"maxNotionalValue"`
	MarginType       string          `json:"marginType"` // "isolated" / "cross" -- lowercase wire format
	IsolatedMargin   decimal.Decimal `json:"isolatedMargin"`
	IsAutoAddMargin  string          `json:"isAutoAddMargin"`
	PositionSide     PositionSide    `json:"positionSide"`
	Notional         decimal.Decimal `json:"notional"`
	IsolatedWallet   decimal.Decimal `json:"isolatedWallet"`
	UpdateTime       time.Time       `json:"updateTime,format:unixmilli"`
}

// GetBuilderUserCommissionRatesService -- GET /fapi/v3/builder/userCommissionRates (USER_DATA)
//
// Queries the commission rates on symbol of the users trading under the
// caller's builder code, computed as GetCommissionRateService does: the user's
// fee tier plus any per-symbol adjustment.
type GetBuilderUserCommissionRatesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserCommissionRatesService(symbol string) *GetBuilderUserCommissionRatesService {
	return &GetBuilderUserCommissionRatesService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetUserAddresses only returns these users (at most 50).
func (s *GetBuilderUserCommissionRatesService) SetUserAddresses(addresses ...string) *GetBuilderUserCommissionRatesService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserCommissionRatesService) SetPage(page int) *GetBuilderUserCommissionRatesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of users per page (default 50, maximum 1000).
func (s *GetBuilderUserCommissionRatesService) SetLimit(limit int) *GetBuilderUserCommissionRatesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserCommissionRatesService) Do(ctx context.Context) (*BuilderUserCommissionRates, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userCommissionRates", s.params).WithSignature()
	return request.Do[BuilderUserCommissionRates](req)
}

// BuilderUserCommissionRates is the paginated response of
// GetBuilderUserCommissionRatesService.
type BuilderUserCommissionRates struct {
	Total       int64                       `json:"total"`
	CurrentPage int                         `json:"currentPage"`
	TotalPages  int                         `json:"totalPages"`
	PageSize    int                         `json:"pageSize"`
	HasMore     bool                        `json:"hasMore"`
	Rows        []BuilderUserCommissionRate `json:"rows"`
	Errors      []BuilderUserError          `json:"errors"`
}

// BuilderUserCommissionRate is a user's commission rates on a symbol.
type BuilderUserCommissionRate struct {
	Address             string          `json:"address"`
	Symbol              string          `json:"symbol"`
	MakerCommissionRate decimal.Decimal `json:"makerCommissionRate"`
	TakerCommissionRate decimal.Decimal `json:"takerCommissionRate"`
}

// GetBuilderUserTradesService -- GET /fapi/v3/builder/userTrades (USER_DATA)
//
// Queries the paginated trade history of users trading under the caller's
// builder code. The authenticated account is the builder identity, so there is
// no separate builder address parameter. If neither startTime nor endTime is
// set, the recent 7 days are returned; the resolved startTime may not be more
// than 90 days before now, and endTime may not be more than 1 day ahead of
// server time. Results are ordered by trade time descending.
type GetBuilderUserTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserTradesService() *GetBuilderUserTradesService {
	return &GetBuilderUserTradesService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns trades of these users (at most 50). Each must
// have approved the caller as its builder; any other address is reported in
// Errors instead. When unset, trades of every bound user are returned.
func (s *GetBuilderUserTradesService) SetUserAddresses(addresses ...string) *GetBuilderUserTradesService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

func (s *GetBuilderUserTradesService) SetStartTime(t time.Time) *GetBuilderUserTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderUserTradesService) SetEndTime(t time.Time) *GetBuilderUserTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderUserTradesService) SetPage(page int) *GetBuilderUserTradesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *GetBuilderUserTradesService) SetLimit(limit int) *GetBuilderUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserTradesService) Do(ctx context.Context) (*BuilderUserTrades, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userTrades", s.params).WithSignature()
	return request.Do[BuilderUserTrades](req)
}

// BuilderUserTrades is the paginated response of GetBuilderUserTradesService.
type BuilderUserTrades struct {
	Total       int64              `json:"total"`
	CurrentPage int                `json:"currentPage"`
	TotalPages  int                `json:"totalPages"`
	PageSize    int                `json:"pageSize"`
	HasMore     bool               `json:"hasMore"`
	Rows        []BuilderUserTrade `json:"rows"`
	Errors      []BuilderUserError `json:"errors"`
}

// BuilderUserTrade is one trade executed by a user under the caller's builder
// code. baseType is CRYPTO or STOCK; totalQuota is the notional (price*qty).
// Fee is from the trading user's side: negative when a fee is charged,
// positive for a rebate.
type BuilderUserTrade struct {
	TradeID        int64           `json:"tradeId"`
	InsertTime     time.Time       `json:"insertTime,format:unixmilli"`
	Symbol         string          `json:"symbol"`
	PositionSide   PositionSide    `json:"positionSide"`
	Price          decimal.Decimal `json:"price"`
	Qty            decimal.Decimal `json:"qty"`
	BaseAsset      string          `json:"baseAsset"`
	BaseType       string          `json:"baseType"`
	QuoteAsset     string          `json:"quoteAsset"`
	Side           OrderSide       `json:"side"`
	ActiveBuy      bool            `json:"activeBuy"`
	FeeAsset       string          `json:"feeAsset"`
	TotalQuota     decimal.Decimal `json:"totalQuota"`
	Fee            decimal.Decimal `json:"fee"`
	OrderID        int64           `json:"orderId"`
	RealizedProfit decimal.Decimal `json:"realizedProfit"`
	MarginAsset    string          `json:"marginAsset"`
	UserAddress    string          `json:"userAddress"`
	BuilderFee     decimal.Decimal `json:"builderFee"`
}

// GetBuilderUserAllOrdersService -- GET /fapi/v3/builder/userAllOrders (USER_DATA)
//
// Queries the paginated order history (active, canceled or filled) of users
// trading under the caller's builder code. If neither startTime nor endTime is
// set, the recent 7 days are returned; the resolved startTime may not be more
// than 90 days before now, and endTime may not be more than 1 day ahead of
// server time.
type GetBuilderUserAllOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserAllOrdersService() *GetBuilderUserAllOrdersService {
	return &GetBuilderUserAllOrdersService{c: c, params: map[string]string{}}
}

// SetUserAddresses only returns orders of these users (at most 50).
func (s *GetBuilderUserAllOrdersService) SetUserAddresses(addresses ...string) *GetBuilderUserAllOrdersService {
	s.params["userAddresses"] = strings.Join(addresses, ",")
	return s
}

// SetSymbol only returns orders on this symbol.
func (s *GetBuilderUserAllOrdersService) SetSymbol(symbol string) *GetBuilderUserAllOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetBuilderUserAllOrdersService) SetStartTime(t time.Time) *GetBuilderUserAllOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderUserAllOrdersService) SetEndTime(t time.Time) *GetBuilderUserAllOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderUserAllOrdersService) SetPage(page int) *GetBuilderUserAllOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of orders per page (default 50, maximum 1000).
func (s *GetBuilderUserAllOrdersService) SetLimit(limit int) *GetBuilderUserAllOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserAllOrdersService) Do(ctx context.Context) (*BuilderUserOrders, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userAllOrders", s.params).WithSignature()
	return request.Do[BuilderUserOrders](req)
}

// BuilderUserOrders is the paginated response of
// GetBuilderUserAllOrdersService.
type BuilderUserOrders struct {
	Total       int64              `json:"total"`
	CurrentPage int                `json:"currentPage"`
	TotalPages  int                `json:"totalPages"`
	PageSize    int                `json:"pageSize"`
	HasMore     bool               `json:"hasMore"`
	Rows        []BuilderUserOrder `json:"rows"`
	Errors      []BuilderUserError `json:"errors"`
}

// BuilderUserOrder is one historical order of a user under the caller's
// builder code. ActivatePrice and PriceRate are only set for
// TRAILING_STOP_MARKET orders, whose StopPrice should be ignored.
type BuilderUserOrder struct {
	Address       string          `json:"address"`
	OrderID       int64           `json:"orderId"`
	Symbol        string          `json:"symbol"`
	Status        OrderStatus     `json:"status"`
	ClientOrderID string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	OrigQty       decimal.Decimal `json:"origQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	CumQuote      decimal.Decimal `json:"cumQuote"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	Type          OrderType       `json:"type"`
	OrigType      OrderType       `json:"origType"`
	ReduceOnly    bool            `json:"reduceOnly"`
	ClosePosition bool            `json:"closePosition"`
	Side          OrderSide       `json:"side"`
	PositionSide  PositionSide    `json:"positionSide"`
	StopPrice     decimal.Decimal `json:"stopPrice"`
	WorkingType   WorkingType     `json:"workingType"`
	PriceProtect  bool            `json:"priceProtect"`
	ActivatePrice decimal.Decimal `json:"activatePrice"`
	PriceRate     decimal.Decimal `json:"priceRate"`
	Time          time.Time       `json:"time,format:unixmilli"`
	UpdateTime    time.Time       `json:"updateTime,format:unixmilli"`
}

// BuilderUserError reports a requested user address that a builder user
// query left out of rows, with the reason: the user has not deposited, has
// privacy mode on, its privacy status could not be checked, or (only when
// userAddresses is set) it is not bound to the caller.
type BuilderUserError struct {
	Address  string `json:"address"`
	ErrorMsg string `json:"errorMsg"`
}

// GetBuilderApprovedUserListService -- GET /fapi/v3/builder/approvedUserList (USER_DATA)
//
// Queries the paginated list of users who have approved the caller's address as
// their builder. The authenticated account is the builder identity; there is no
// separate builder address parameter. Results are returned page by page under
// rows[]; when no time filter is set the server returns all approved users.
type GetBuilderApprovedUserListService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderApprovedUserListService() *GetBuilderApprovedUserListService {
	return &GetBuilderApprovedUserListService{c: c, params: map[string]string{}}
}

// SetStartTime only returns users approved at or after this timestamp. When
// unset the server applies no lower time bound.
func (s *GetBuilderApprovedUserListService) SetStartTime(t time.Time) *GetBuilderApprovedUserListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime only returns users approved at or before this timestamp. When
// unset the server applies no upper time bound.
func (s *GetBuilderApprovedUserListService) SetEndTime(t time.Time) *GetBuilderApprovedUserListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPage sets the pagination page number (default 1).
func (s *GetBuilderApprovedUserListService) SetPage(page int) *GetBuilderApprovedUserListService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of records per page (default 100, maximum 1000).
func (s *GetBuilderApprovedUserListService) SetLimit(limit int) *GetBuilderApprovedUserListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderApprovedUserListService) Do(ctx context.Context) (*BuilderApprovedUserList, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/approvedUserList", s.params).WithSignature()
	return request.Do[BuilderApprovedUserList](req)
}

// BuilderApprovedUserList is the paginated response of
// GetBuilderApprovedUserListService.
type BuilderApprovedUserList struct {
	Total       int64                 `json:"total"`
	CurrentPage int                   `json:"currentPage"`
	TotalPages  int                   `json:"totalPages"`
	PageSize    int                   `json:"pageSize"`
	HasMore     bool                  `json:"hasMore"`
	Rows        []BuilderApprovedUser `json:"rows"`
}

// BuilderApprovedUser is one user who approved the caller as their builder.
type BuilderApprovedUser struct {
	UserAddress    string          `json:"userAddress"`
	BuilderAddress string          `json:"builderAddress"`
	MaxFeeRate     decimal.Decimal `json:"maxFeeRate"`
	BuilderName    string          `json:"builderName"`
	ApproveTime    time.Time       `json:"approveTime,format:unixmilli"`
}
