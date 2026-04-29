package futures

// OrderSide -- BUY or SELL.
type OrderSide string

const (
	SideBuy  OrderSide = "BUY"
	SideSell OrderSide = "SELL"
)

// PositionSide is BOTH for one-way mode, LONG/SHORT for hedge mode.
type PositionSide string

const (
	PositionSideBoth  PositionSide = "BOTH"
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

// OrderType covers all V3 futures order types (TRAILING_STOP_MARKET extra
// over spot).
type OrderType string

const (
	OrderTypeLimit              OrderType = "LIMIT"
	OrderTypeMarket             OrderType = "MARKET"
	OrderTypeStop               OrderType = "STOP"
	OrderTypeStopMarket         OrderType = "STOP_MARKET"
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET"
)

// OrderStatus mirrors the matching engine.
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce values accepted by the futures API.
type TimeInForce string

const (
	TimeInForceGTC    TimeInForce = "GTC"
	TimeInForceIOC    TimeInForce = "IOC"
	TimeInForceFOK    TimeInForce = "FOK"
	TimeInForceGTX    TimeInForce = "GTX"
	TimeInForceHidden TimeInForce = "HIDDEN"
)

// WorkingType determines what price triggers stop/take-profit orders.
type WorkingType string

const (
	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
)

// MarginType -- ISOLATED or CROSSED.
type MarginType string

const (
	MarginTypeIsolated MarginType = "ISOLATED"
	MarginTypeCrossed  MarginType = "CROSSED"
)

// ResponseType controls how detailed POST /order response is.
type ResponseType string

const (
	ResponseTypeAck    ResponseType = "ACK"
	ResponseTypeResult ResponseType = "RESULT"
)

// PositionMarginType discriminates between adding (1) and removing (2) margin
// on isolated positions in ModifyIsolatedPositionMargin.
type PositionMarginType int

const (
	PositionMarginAdd    PositionMarginType = 1
	PositionMarginReduce PositionMarginType = 2
)

// AutoCloseType filters force-orders by reason.
type AutoCloseType string

const (
	AutoCloseLiquidation AutoCloseType = "LIQUIDATION"
	AutoCloseADL         AutoCloseType = "ADL"
)

// IncomeType filters /fapi/v3/income.
type IncomeType string

const (
	IncomeTransfer                IncomeType = "TRANSFER"
	IncomeWelcomeBonus            IncomeType = "WELCOME_BONUS"
	IncomeRealizedPNL             IncomeType = "REALIZED_PNL"
	IncomeFundingFee              IncomeType = "FUNDING_FEE"
	IncomeCommission              IncomeType = "COMMISSION"
	IncomeInsuranceClear          IncomeType = "INSURANCE_CLEAR"
	IncomeMarketMerchantReturnRwd IncomeType = "MARKET_MERCHANT_RETURN_REWARD"
)

// TransferKindType controls direction of a futures<->spot transfer.
type TransferKindType string

const (
	TransferFutureToSpot TransferKindType = "FUTURE_SPOT"
	TransferSpotToFuture TransferKindType = "SPOT_FUTURE"
)

// KlineInterval values accepted by /fapi/v3/klines.
type KlineInterval string

const (
	Interval1m  KlineInterval = "1m"
	Interval3m  KlineInterval = "3m"
	Interval5m  KlineInterval = "5m"
	Interval15m KlineInterval = "15m"
	Interval30m KlineInterval = "30m"
	Interval1h  KlineInterval = "1h"
	Interval2h  KlineInterval = "2h"
	Interval4h  KlineInterval = "4h"
	Interval6h  KlineInterval = "6h"
	Interval8h  KlineInterval = "8h"
	Interval12h KlineInterval = "12h"
	Interval1d  KlineInterval = "1d"
	Interval3d  KlineInterval = "3d"
	Interval1w  KlineInterval = "1w"
	Interval1M  KlineInterval = "1M"
)
