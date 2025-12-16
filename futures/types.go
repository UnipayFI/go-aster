package futures

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type PositionSide string

const (
	PositionSideBoth  PositionSide = "BOTH"
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

type MarginType string

const (
	MarginTypeIsolated MarginType = "ISOLATED"
	MarginTypeCrossed  MarginType = "CROSSED"
)

type WorkingType string

const (
	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
)

type IncomeType string

const (
	IncomeTypeTransfer       IncomeType = "TRANSFER"
	IncomeTypeWelcomeBonus   IncomeType = "WELCOME_BONUS"
	IncomeTypeRealizedPnl    IncomeType = "REALIZED_PNL"
	IncomeTypeFundingFee     IncomeType = "FUNDING_FEE"
	IncomeTypeCommission     IncomeType = "COMMISSION"
	IncomeTypeInsuranceClear IncomeType = "INSURANCE_CLEAR"
	IncomeTypeMarketMerchant IncomeType = "MARKET_MERCHANT_RETURN_REWARD"
)

type ContractType string

const (
	ContractTypePerpetual ContractType = "PERPETUAL"
)

type ContractStatus string

const (
	ContractStatusPendingTrading ContractStatus = "PENDING_TRADING"
	ContractStatusTrading        ContractStatus = "TRADING"
	ContractStatusPreSettle      ContractStatus = "PRE_SETTLE"
	ContractStatusSettling       ContractStatus = "SETTLING"
	ContractStatusClose          ContractStatus = "CLOSE"
)

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

type TimeInForce string

const (
	TimeInForceGTC    TimeInForce = "GTC"
	TimeInForceIOC    TimeInForce = "IOC"
	TimeInForceFOK    TimeInForce = "FOK"
	TimeInForceGTX    TimeInForce = "GTX"
	TimeInForceHidden TimeInForce = "HIDDEN"
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

type NewOrderRespType string

const (
	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeResult NewOrderRespType = "RESULT"
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

type TransferType string

const (
	TransferTypeFutureToSpot TransferType = "FUTURE_SPOT"
	TransferTypeSpotToFuture TransferType = "SPOT_FUTURE"
)

type AutoCloseType string

const (
	AutoCloseTypeLiquidation AutoCloseType = "LIQUIDATION"
	AutoCloseTypeADL         AutoCloseType = "ADL"
)

type DepthLevel struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}

type GeneralResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handlerGeneralResponse(resp *GeneralResponse, err error) error {
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return fmt.Errorf("%s failed: %d %v", resp.Message, resp.Code, resp.Message)
	}
	return nil
}
