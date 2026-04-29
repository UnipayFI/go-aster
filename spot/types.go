package spot

// OrderSide -- BUY or SELL.
type OrderSide string

const (
	SideBuy  OrderSide = "BUY"
	SideSell OrderSide = "SELL"
)

// OrderType covers all V3 spot order types.
type OrderType string

const (
	OrderTypeLimit             OrderType = "LIMIT"
	OrderTypeMarket            OrderType = "MARKET"
	OrderTypeStop              OrderType = "STOP"
	OrderTypeStopMarket        OrderType = "STOP_MARKET"
	OrderTypeTakeProfit        OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket  OrderType = "TAKE_PROFIT_MARKET"
)

// OrderStatus from the matching engine.
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce determines how long an order stays active.
type TimeInForce string

const (
	TimeInForceGTC    TimeInForce = "GTC"
	TimeInForceIOC    TimeInForce = "IOC"
	TimeInForceFOK    TimeInForce = "FOK"
	TimeInForceGTX    TimeInForce = "GTX"
	TimeInForceHidden TimeInForce = "HIDDEN"
)

// KlineInterval values accepted by /api/v3/klines.
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

// TransferKindType controls direction of the perp/spot internal transfer.
type TransferKindType string

const (
	TransferFutureToSpot TransferKindType = "FUTURE_SPOT"
	TransferSpotToFuture TransferKindType = "SPOT_FUTURE"
)

// TransactionType filters returned by /api/v3/transactionHistory.
type TransactionType string

const (
	TransactionTradeTarget          TransactionType = "TRADE_TARGET"
	TransactionTradeSource          TransactionType = "TRADE_SOURCE"
	TransactionTransferSpotToFuture TransactionType = "TRANSFER_SPOT_TO_FUTURE"
	TransactionTransferFutureToSpot TransactionType = "TRANSFER_FUTURE_TO_SPOT"
	TransactionTransferSpotToSpot   TransactionType = "TRANSFER_SPOT_TO_SPOT"
	TransactionAirdrop              TransactionType = "AIRDROP"
	TransactionDividend             TransactionType = "DIVIDEND"
	TransactionTransferRefund       TransactionType = "TRANSFER_REFUND"
	TransactionInternalTransfer     TransactionType = "INTERNAL_TRANSFER"
	TransactionTransfer             TransactionType = "TRANSFER"
	TransactionSwap                 TransactionType = "SWAP"
	TransactionCommissionRebate     TransactionType = "COMMISSION_REBATE"
	TransactionCashBack             TransactionType = "CASH_BACK"
	TransactionStakingWithdraw      TransactionType = "STAKING_WITHDRAW"
	TransactionStakingClaim         TransactionType = "STAKING_CLAIM"
	TransactionStakingDelegate      TransactionType = "STAKING_DELEGATE"
)
