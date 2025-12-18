package futures

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/request"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

type SubscribeUserDataStreamService struct {
	c         *FuturesWebSocketClient
	listenKey string
}

func (c *FuturesWebSocketClient) NewSubscribeUserDataStreamService(listenKey string) *SubscribeUserDataStreamService {
	return &SubscribeUserDataStreamService{c: c, listenKey: listenKey}
}

func (s *SubscribeUserDataStreamService) Do(ctx context.Context, handler UserDataHandler) (doneC <-chan struct{}, stopC chan struct{}, err error) {
	callback := func(message *json.RawMessage, err error) {
		if err != nil {
			handler.OnError(err)
			return
		}

		switch UserDataEventType(gjson.GetBytes(*message, "e").String()) {
		case AccountUpdate:
			handleEvent(s, message, &WsAccountUpdateEvent{}, handler.OnAccountUpdate, handler.OnError)
		case OrderTradeUpdate:
			handleEvent(s, message, &WsOrderTradeUpdateEvent{}, handler.OnOrderTradeUpdate, handler.OnError)
		case AccountConfigUpdate:
			handleEvent(s, message, &WsAccountConfigUpdateEvent{}, handler.OnAccountConfigUpdate, handler.OnError)
		case ListenKeyExpired:
			handleEvent(s, message, &WsListenKeyExpiredEvent{}, handler.OnListenKeyExpired, handler.OnError)
		default:
			handler.OnError(fmt.Errorf("unknown event type: %s", gjson.GetBytes(*message, "e").String()))
			return
		}
	}

	return request.Subscribe(ctx, s.c, common.WEBSOCKET_STREAM_SEPARATOR+s.listenKey, callback)
}

func handleEvent[T any](s *SubscribeUserDataStreamService, message *json.RawMessage, target *T, handle func(*T), onError func(error)) {
	if err := s.c.GetHttpClient().JSONUnmarshal(*message, target); err != nil {
		onError(err)
		return
	}
	handle(target)
}

type UserDataHandler interface {
	OnAccountUpdate(message *WsAccountUpdateEvent)
	OnOrderTradeUpdate(message *WsOrderTradeUpdateEvent)
	OnAccountConfigUpdate(message *WsAccountConfigUpdateEvent)
	OnListenKeyExpired(message *WsListenKeyExpiredEvent)
	OnError(err error)
}

type WsAccountUpdateEvent struct {
	WsBaseEvent
	TransactionTime time.Time       `json:"T,format:unixmilli"`
	AccountUpdate   WsAccountUpdate `json:"a"`
}

type WsAccountUpdate struct {
	Reason    AccountUpdateReason `json:"m"`
	Balances  []WsAccountBalance  `json:"B"`
	Positions []WsAccountPosition `json:"P"`
}

type WsAccountBalance struct {
	Asset              string          `json:"a"`
	WalletBalance      decimal.Decimal `json:"wb"`
	CrossWalletBalance decimal.Decimal `json:"cw"`
	ChangeBalance      decimal.Decimal `json:"bc"`
}

type WsAccountPosition struct {
	Symbol              string          `json:"s"`
	PositionAmount      decimal.Decimal `json:"pa"`
	EntryPrice          decimal.Decimal `json:"ep"`
	AccumulatedRealized decimal.Decimal `json:"cr"`
	UnrealizedPnL       decimal.Decimal `json:"up"`
	MarginType          MarginType      `json:"mt"`
	IsolatedWallet      decimal.Decimal `json:"iw"`
	Side                PositionSide    `json:"ps"`
}

type WsOrderTradeUpdateEvent struct {
	WsBaseEvent
	TransactionTime  time.Time          `json:"T,format:unixmilli"`
	OrderTradeUpdate WsOrderTradeUpdate `json:"o"`
}

type WsOrderTradeUpdate struct {
	Symbol               string          `json:"s"`                  // Symbol
	ClientOrderId        string          `json:"c"`                  // Client order ID
	Side                 OrderSide       `json:"S"`                  // Side
	OrderType            OrderType       `json:"o"`                  // Order type
	TimeInForce          TimeInForce     `json:"f"`                  // Time in force
	OriginalQty          decimal.Decimal `json:"q"`                  // Original quantity
	OriginalPrice        decimal.Decimal `json:"p"`                  // Original price
	AveragePrice         decimal.Decimal `json:"ap"`                 // Average price
	StopPrice            decimal.Decimal `json:"sp"`                 // Stop price. Please ignore with TRAILING_STOP_MARKET order
	ExecutionType        ExecutionType   `json:"x"`                  // Execution type
	Status               OrderStatus     `json:"X"`                  // Order status
	ID                   int64           `json:"i"`                  // Order ID
	LastFilledQty        decimal.Decimal `json:"l"`                  // Order Last Filled Quantity
	AccumulatedFilledQty decimal.Decimal `json:"z"`                  // Order Filled Accumulated Quantity
	LastFilledPrice      decimal.Decimal `json:"L"`                  // Last Filled Price
	CommissionAsset      string          `json:"N"`                  // Commission Asset, will not push if no commission
	CommissionAmount     decimal.Decimal `json:"n"`                  // Commission, will not push if no commission
	TradeTime            time.Time       `json:"T,format:unixmilli"` // Order Trade Time
	TradeID              int64           `json:"t"`                  // Trade ID
	BidsNotional         decimal.Decimal `json:"b"`                  // Bids Notional
	AsksNotional         decimal.Decimal `json:"a"`                  // Asks Notional
	IsMaker              bool            `json:"m"`                  // Is this trade the maker side?
	IsReduceOnly         bool            `json:"R"`                  // Is this reduce only
	WorkingType          WorkingType     `json:"wt"`                 // Stop Price Working Type
	OriginalOrderType    OrderType       `json:"ot"`                 // Original Order Type
	PositionSide         PositionSide    `json:"ps"`                 // Position Side
	IsClosingPosition    bool            `json:"cp"`                 // If Close-All, pushed with conditional order
	ActivationPrice      decimal.Decimal `json:"AP"`                 // Activation Price, only puhed with TRAILING_STOP_MARKET order
	CallbackRate         decimal.Decimal `json:"cr"`                 //Callback Rate, only puhed with TRAILING_STOP_MARKET order
	RealizedProfit       decimal.Decimal `json:"rp"`                 // Realized Profit of the trade
}

type ExecutionType string

const (
	NEW        ExecutionType = "NEW"
	CANCELED   ExecutionType = "CANCELED"
	CALCULATED ExecutionType = "CALCULATED"
	EXPIRED    ExecutionType = "EXPIRED"
	TRADE      ExecutionType = "TRADE"
)

type WsAccountConfigUpdateEvent struct {
	WsBaseEvent
	TransactionTime      time.Time               `json:"T,format:unixmilli"`
	SymbolLeverageUpdate *WsSymbolLeverageUpdate `json:"ac"`
	UserAccountUpdate    *WsUserAccountUpdate    `json:"u"`
}

type WsSymbolLeverageUpdate struct {
	Symbol   string `json:"s"`
	Leverage int64  `json:"l"`
}

type WsUserAccountUpdate struct {
	MultiAssetsMode            bool `json:"j"` // Multi-Assets Mode
	SpecifiedTokenFeeDeduction bool `json:"f"` // Specified token fee deduction
	PositionMode               bool `json:"d"` // Position mode: true for dual-side (hedge) mode, false for single-side (one-way) mode
}

type WsListenKeyExpiredEvent struct {
	WsBaseEvent
}

type WsBaseEvent struct {
	Event UserDataEventType `json:"e"`
	Time  time.Time         `json:"E,format:unixmilli"`
}

type AccountUpdateReason string

const (
	DEPOSIT               AccountUpdateReason = "DEPOSIT"
	WITHDRAW              AccountUpdateReason = "WITHDRAW"
	ORDER                 AccountUpdateReason = "ORDER"
	FUNDING_FEE           AccountUpdateReason = "FUNDING_FEE"
	WITHDRAW_REJECT       AccountUpdateReason = "WITHDRAW_REJECT"
	ADJUSTMENT            AccountUpdateReason = "ADJUSTMENT"
	INSURANCE_CLEAR       AccountUpdateReason = "INSURANCE_CLEAR"
	ADMIN_DEPOSIT         AccountUpdateReason = "ADMIN_DEPOSIT"
	ADMIN_WITHDRAW        AccountUpdateReason = "ADMIN_WITHDRAW"
	MARGIN_TRANSFER       AccountUpdateReason = "MARGIN_TRANSFER"
	MARGIN_TYPE_CHANGE    AccountUpdateReason = "MARGIN_TYPE_CHANGE"
	ASSET_TRANSFER        AccountUpdateReason = "ASSET_TRANSFER"
	OPTIONS_PREMIUM_FEE   AccountUpdateReason = "OPTIONS_PREMIUM_FEE"
	OPTIONS_SETTLE_PROFIT AccountUpdateReason = "OPTIONS_SETTLE_PROFIT"
	AUTO_EXCHANGE         AccountUpdateReason = "AUTO_EXCHANGE"
)

type UserDataEventType string

const (
	AccountUpdate       UserDataEventType = "ACCOUNT_UPDATE"
	OrderTradeUpdate    UserDataEventType = "ORDER_TRADE_UPDATE"
	AccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE"
	ListenKeyExpired    UserDataEventType = "listenKeyExpired"
)
