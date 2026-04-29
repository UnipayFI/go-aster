package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

// CreateListenKeyService -- POST /api/v3/listenKey (USER_STREAM)
type CreateListenKeyService struct {
	c *SpotClient
}

func (c *SpotClient) NewCreateListenKeyService() *CreateListenKeyService {
	return &CreateListenKeyService{c: c}
}

func (s *CreateListenKeyService) Do(ctx context.Context) (*ListenKeyResponse, error) {
	req := request.Post(s.c, ctx, "/api/v3/listenKey").WithSignature()
	return request.Do[ListenKeyResponse](req)
}

type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

// RenewListenKeyService -- PUT /api/v3/listenKey (USER_STREAM)
type RenewListenKeyService struct {
	c         *SpotClient
	listenKey string
}

func (c *SpotClient) NewRenewListenKeyService(listenKey string) *RenewListenKeyService {
	return &RenewListenKeyService{c: c, listenKey: listenKey}
}

func (s *RenewListenKeyService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/api/v3/listenKey", map[string]string{"listenKey": s.listenKey}).WithSignature()
	_, err := request.Do[struct{}](req)
	return err
}

// DeleteListenKeyService -- DELETE /api/v3/listenKey (USER_STREAM)
type DeleteListenKeyService struct {
	c         *SpotClient
	listenKey string
}

func (c *SpotClient) NewDeleteListenKeyService(listenKey string) *DeleteListenKeyService {
	return &DeleteListenKeyService{c: c, listenKey: listenKey}
}

func (s *DeleteListenKeyService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v3/listenKey", map[string]string{"listenKey": s.listenKey}).WithSignature()
	_, err := request.Do[struct{}](req)
	return err
}

// SubscribeUserDataStreamService streams the discriminated union of events
// the listenKey-based stream emits: balance changes, execution reports, and
// tradepro fills. The callback receives a fully-decoded WsUserDataEvent
// where exactly one of AccountUpdate / ExecutionReport / TradePro is non-nil
// based on the "e" field. Unknown event types still arrive with EventType
// set and the original bytes available in Raw for forward compatibility.
type SubscribeUserDataStreamService struct {
	c         *SpotWebSocketClient
	listenKey string
}

func (c *SpotWebSocketClient) NewSubscribeUserDataStreamService(listenKey string) *SubscribeUserDataStreamService {
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
	case "outboundAccountPosition":
		var x WsAccountUpdateEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.AccountUpdate = &x
		ev.EventTime = x.EventTime
	case "executionReport":
		var x WsExecutionReportEvent
		if err := json.Unmarshal(msg, &x); err != nil {
			return ev, err
		}
		ev.ExecutionReport = &x
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

// WsUserDataEvent is a discriminated union. Inspect EventType, then read
// from the matching pointer field; the others remain nil.
type WsUserDataEvent struct {
	EventType       string
	EventTime       time.Time
	Raw             []byte
	AccountUpdate   *WsAccountUpdateEvent
	ExecutionReport *WsExecutionReportEvent
	TradePro        *WsTradeProEvent
}

type WsAccountUpdateEvent struct {
	EventType      string      `json:"e"`
	EventTime      time.Time   `json:"E,format:unixmilli"`
	LastUpdateTime time.Time   `json:"T,format:unixmilli"`
	Balances       []WsBalance `json:"B"`
	Reason         string      `json:"m"`
}

type WsBalance struct {
	Asset  string          `json:"a"`
	Free   decimal.Decimal `json:"f"`
	Locked decimal.Decimal `json:"l"`
}

// WsExecutionReportEvent corresponds to e=="executionReport".
//
// ExecutionType values: NEW / CANCELED / REJECTED / TRADE / EXPIRED.
type WsExecutionReportEvent struct {
	EventType                string          `json:"e"`
	EventTime                time.Time       `json:"E,format:unixmilli"`
	Symbol                   string          `json:"s"`
	ClientOrderID            string          `json:"c"`
	Side                     OrderSide       `json:"S"`
	OrderType                OrderType       `json:"o"`
	TimeInForce              TimeInForce     `json:"f"`
	OrderQuantity            decimal.Decimal `json:"q"`
	OrderPrice               decimal.Decimal `json:"p"`
	AvgPrice                 decimal.Decimal `json:"ap"`
	StopPrice                decimal.Decimal `json:"P"`
	ExecutionType            string          `json:"x"`
	OrderStatus              OrderStatus     `json:"X"`
	OrderID                  int64           `json:"i"`
	LastFilledQty            decimal.Decimal `json:"l"`
	CumFilledQty             decimal.Decimal `json:"z"`
	LastFilledPrice          decimal.Decimal `json:"L"`
	CommissionAmount         decimal.Decimal `json:"n"`
	CommissionAsset          string          `json:"N"`
	TransactionTime          time.Time       `json:"T,format:unixmilli"`
	TradeID                  int64           `json:"t"`
	IsMaker                  bool            `json:"m"`
	OrigOrderType            OrderType       `json:"ot"`
	OrderCreationTime        time.Time       `json:"O,format:unixmilli"`
	CumQuoteAssetTransacted  decimal.Decimal `json:"Z"`
	LastQuoteAssetTransacted decimal.Decimal `json:"Y"`
	QuoteOrderQty            decimal.Decimal `json:"Q"`
}
