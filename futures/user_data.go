package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

// CreateListenKeyService -- POST /fapi/v3/listenKey (USER_STREAM)
type CreateListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewCreateListenKeyService() *CreateListenKeyService {
	return &CreateListenKeyService{c: c}
}

func (s *CreateListenKeyService) Do(ctx context.Context) (*ListenKeyResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/listenKey").WithSignature()
	return request.Do[ListenKeyResponse](req)
}

type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

// RenewListenKeyService -- PUT /fapi/v3/listenKey (USER_STREAM)
//
// Unlike spot, futures keepalive takes no params (the listenKey is implied
// by the signer). Recommended to call every ~60 minutes.
type RenewListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewRenewListenKeyService() *RenewListenKeyService {
	return &RenewListenKeyService{c: c}
}

func (s *RenewListenKeyService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/fapi/v3/listenKey").WithSignature()
	_, err := request.Do[struct{}](req)
	return err
}

// DeleteListenKeyService -- DELETE /fapi/v3/listenKey (USER_STREAM)
type DeleteListenKeyService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewDeleteListenKeyService() *DeleteListenKeyService {
	return &DeleteListenKeyService{c: c}
}

func (s *DeleteListenKeyService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/fapi/v3/listenKey").WithSignature()
	_, err := request.Do[struct{}](req)
	return err
}

// SubscribeUserDataStreamService streams the discriminated union of futures
// user-data events. Inspect EventType, then read the matching pointer field.
type SubscribeUserDataStreamService struct {
	c         *FuturesWebSocketClient
	listenKey string
}

func (c *FuturesWebSocketClient) NewSubscribeUserDataStreamService(listenKey string) *SubscribeUserDataStreamService {
	return &SubscribeUserDataStreamService{c: c, listenKey: listenKey}
}

func (s *SubscribeUserDataStreamService) Do(ctx context.Context, cb func(*WsUserDataEvent, error)) (chan<- struct{}, <-chan struct{}, error) {
	endpoint := common.WEBSOCKET_STREAM_SEPARATOR + s.listenKey
	return request.SubscribeRaw(ctx, s.c, endpoint, func(msg []byte, err error) {
		if err != nil {
			cb(nil, err)
			return
		}
		ev, decodeErr := decodeUserDataEvent(msg)
		cb(ev, decodeErr)
	})
}

func decodeUserDataEvent(msg []byte) (*WsUserDataEvent, error) {
	var head struct {
		EventType string `json:"e"`
	}
	if err := json.Unmarshal(msg, &head); err != nil {
		return nil, err
	}
	ev := &WsUserDataEvent{
		EventType: head.EventType,
		Raw:       msg,
	}
	switch head.EventType {
	case "listenKeyExpired":
		var x WsListenKeyExpiredEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.ListenKeyExpired = &x
		ev.EventTime = x.EventTime
	case "MARGIN_CALL":
		var x WsMarginCallEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.MarginCall = &x
		ev.EventTime = x.EventTime
	case "ACCOUNT_UPDATE":
		var x WsAccountUpdateEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.AccountUpdate = &x
		ev.EventTime = x.EventTime
	case "ORDER_TRADE_UPDATE":
		var x WsOrderTradeUpdateEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.OrderTradeUpdate = &x
		ev.EventTime = x.EventTime
	case "ACCOUNT_CONFIG_UPDATE":
		var x WsAccountConfigUpdateEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.AccountConfigUpdate = &x
		ev.EventTime = x.EventTime
	case "tradepro":
		var x WsTradeProEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.TradePro = &x
		ev.EventTime = x.EventTime
	}
	return ev, nil
}

// WsUserDataEvent is the discriminated union of every event the futures
// user-data stream can deliver. Inspect EventType to know which pointer
// field is populated; Raw is always present for forward-compat.
type WsUserDataEvent struct {
	EventType           string
	EventTime           time.Time
	Raw                 []byte
	ListenKeyExpired    *WsListenKeyExpiredEvent
	MarginCall          *WsMarginCallEvent
	AccountUpdate       *WsAccountUpdateEvent
	OrderTradeUpdate    *WsOrderTradeUpdateEvent
	AccountConfigUpdate *WsAccountConfigUpdateEvent
	TradePro            *WsTradeProEvent
}

type WsListenKeyExpiredEvent struct {
	EventType string    `json:"e"`
	EventTime time.Time `json:"E,format:unixmilli"`
}

type WsMarginCallEvent struct {
	EventType          string                  `json:"e"`
	EventTime          time.Time               `json:"E,format:unixmilli"`
	CrossWalletBalance decimal.Decimal         `json:"cw"`
	Positions          []WsMarginCallPosition  `json:"p"`
}

type WsMarginCallPosition struct {
	Symbol         string          `json:"s"`
	PositionSide   PositionSide    `json:"ps"`
	PositionAmount decimal.Decimal `json:"pa"`
	MarginType     string          `json:"mt"`
	IsolatedWallet decimal.Decimal `json:"iw"`
	MarkPrice      decimal.Decimal `json:"mp"`
	UnrealizedPnL  decimal.Decimal `json:"up"`
	MaintMargin    decimal.Decimal `json:"mm"`
}

type WsAccountUpdateEvent struct {
	EventType       string             `json:"e"`
	EventTime       time.Time          `json:"E,format:unixmilli"`
	TransactionTime time.Time          `json:"T,format:unixmilli"`
	UpdateData      WsAccountUpdateRO  `json:"a"`
}

type WsAccountUpdateRO struct {
	EventReasonType string                  `json:"m"`
	Balances        []WsAccountUpdateBal    `json:"B"`
	Positions       []WsAccountUpdatePos    `json:"P"`
}

type WsAccountUpdateBal struct {
	Asset              string          `json:"a"`
	WalletBalance      decimal.Decimal `json:"wb"`
	CrossWalletBalance decimal.Decimal `json:"cw"`
	BalanceChange      decimal.Decimal `json:"bc"`
}

type WsAccountUpdatePos struct {
	Symbol               string          `json:"s"`
	PositionAmount       decimal.Decimal `json:"pa"`
	EntryPrice           decimal.Decimal `json:"ep"`
	AccumulatedRealized  decimal.Decimal `json:"cr"`
	UnrealizedPnL        decimal.Decimal `json:"up"`
	MarginType           string          `json:"mt"`
	IsolatedWallet       decimal.Decimal `json:"iw"`
	PositionSide         PositionSide    `json:"ps"`
}

// WsOrderTradeUpdateEvent corresponds to e=="ORDER_TRADE_UPDATE".
type WsOrderTradeUpdateEvent struct {
	EventType       string                 `json:"e"`
	EventTime       time.Time              `json:"E,format:unixmilli"`
	TransactionTime time.Time              `json:"T,format:unixmilli"`
	Order           WsOrderTradeUpdateRO   `json:"o"`
}

// WsOrderTradeUpdateRO holds the order-update payload. ExecutionType values:
// NEW / CANCELED / CALCULATED / EXPIRED / TRADE. OrderStatus may also include
// NEW_INSURANCE / NEW_ADL for liquidation paths.
type WsOrderTradeUpdateRO struct {
	Symbol                 string          `json:"s"`
	ClientOrderID          string          `json:"c"`
	Side                   OrderSide       `json:"S"`
	OrderType              OrderType       `json:"o"`
	TimeInForce            TimeInForce     `json:"f"`
	OrigQty                decimal.Decimal `json:"q"`
	OrigPrice              decimal.Decimal `json:"p"`
	AvgPrice               decimal.Decimal `json:"ap"`
	StopPrice              decimal.Decimal `json:"sp"`
	ExecutionType          string          `json:"x"`
	OrderStatus            string          `json:"X"`
	OrderID                int64           `json:"i"`
	LastFilledQty          decimal.Decimal `json:"l"`
	AccumFilledQty         decimal.Decimal `json:"z"`
	LastFilledPrice        decimal.Decimal `json:"L"`
	CommissionAsset        string          `json:"N"`
	Commission             decimal.Decimal `json:"n"`
	OrderTradeTime         time.Time       `json:"T,format:unixmilli"`
	TradeID                int64           `json:"t"`
	BidsNotional           decimal.Decimal `json:"b"`
	AskNotional            decimal.Decimal `json:"a"`
	IsMaker                bool            `json:"m"`
	IsReduceOnly           bool            `json:"R"`
	StopPriceWorkingType   WorkingType     `json:"wt"`
	OrigOrderType          OrderType       `json:"ot"`
	PositionSide           PositionSide    `json:"ps"`
	ClosePosition          bool            `json:"cp"`
	ActivationPrice        decimal.Decimal `json:"AP"`
	CallbackRate           decimal.Decimal `json:"cr"`
	RealizedProfit         decimal.Decimal `json:"rp"`
}

// WsAccountConfigUpdateEvent corresponds to e=="ACCOUNT_CONFIG_UPDATE".
//
// One of TradingPair / AccountInfo will be populated; the other stays at its
// zero value. Inspect TradingPair.Symbol or AccountInfo (zero check) to tell.
type WsAccountConfigUpdateEvent struct {
	EventType       string                       `json:"e"`
	EventTime       time.Time                    `json:"E,format:unixmilli"`
	TransactionTime time.Time                    `json:"T,format:unixmilli"`
	TradingPair     *WsAccountConfigTradingPair  `json:"ac,omitempty"`
	AccountInfo     *WsAccountConfigAccountInfo  `json:"ai,omitempty"`
}

type WsAccountConfigTradingPair struct {
	Symbol   string          `json:"s"`
	Leverage decimal.Decimal `json:"l"`
}

type WsAccountConfigAccountInfo struct {
	MultiAssetsMode    bool `json:"j"`
	TokenFeeDeduction  bool `json:"f"`
	HedgeMode          bool `json:"d"`
}

// WsTradeProEvent corresponds to e=="tradepro" (Aster on-chain trade event).
//
// Note futures tradepro arrives as a combined-stream wrapped payload (with a
// leading "stream" envelope), whereas the listenKey-based user stream
// delivers it as a flat document. We accept the flat form here; if you
// subscribe via /stream?streams=, unwrap "data" first.
type WsTradeProEvent struct {
	EventType       string          `json:"e"`
	EventTime       time.Time       `json:"E,format:unixmilli"`
	TradeTime       time.Time       `json:"T,format:unixmilli"`
	Symbol          string          `json:"s"`
	TradeID         int64           `json:"t"`
	Price           decimal.Decimal `json:"p"`
	Quantity        decimal.Decimal `json:"q"`
	TransactionHash string          `json:"h"`
	Participants    []string        `json:"m"`
}
