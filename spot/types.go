package spot

type OrderType string

const (
	OrderTypeLimit            OrderType = "LIMIT"
	OrderTypeMarket           OrderType = "MARKET"
	OrderTypeStop             OrderType = "STOP"
	OrderTypeStopMarket       OrderType = "STOP_MARKET"
	OrderTypeTakeProfit       OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket OrderType = "TAKE_PROFIT_MARKET"
)

type TimeInForce string

const (
	TimeInForceTypeGTC    TimeInForce = "GTC"    // Good Till Cancel
	TimeInForceTypeIOC    TimeInForce = "IOC"    // Immediate or Cancel
	TimeInForceTypeFOK    TimeInForce = "FOK"    // Fill or Kill
	TimeInForceTypeGTX    TimeInForce = "GTX"    // Good Till Crossing (Post Only)
	TimeInForceTypeHIDDEN TimeInForce = "HIDDEN" // Hidden Order (Post Only)
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

type TransferType string

const (
	TransferFutureToSpot TransferType = "FUTURE_SPOT"
	TransferSpotToFuture TransferType = "SPOT_FUTURE"
)

type KlineInterval string

const (
	KlineInterval1m  KlineInterval = "1m"
	KlineInterval3m  KlineInterval = "3m"
	KlineInterval5m  KlineInterval = "5m"
	KlineInterval15m KlineInterval = "15m"
	KlineInterval30m KlineInterval = "30m"
	KlineInterval1h  KlineInterval = "1h"
	KlineInterval2h  KlineInterval = "2h"
	KlineInterval4h  KlineInterval = "4h"
	KlineInterval6h  KlineInterval = "6h"
	KlineInterval8h  KlineInterval = "8h"
	KlineInterval12h KlineInterval = "12h"
	KlineInterval1d  KlineInterval = "1d"
	KlineInterval3d  KlineInterval = "3d"
	KlineInterval1w  KlineInterval = "1w"
	KlineInterval1M  KlineInterval = "1M"
)
