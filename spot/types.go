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
