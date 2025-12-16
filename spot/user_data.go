package spot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

type SubscribeUserDataStreamService struct {
	c         *SpotWebSocketClient
	listenKey string
}

func (c *SpotWebSocketClient) NewSubscribeUserDataStreamService(listenKey string) *SubscribeUserDataStreamService {
	return &SubscribeUserDataStreamService{c: c, listenKey: listenKey}
}

func (s *SubscribeUserDataStreamService) Do(ctx context.Context, handler UserDataHandler) (doneC <-chan struct{}, stopC chan struct{}, err error) {
	callback := func(message *json.RawMessage, err error) {
		if err != nil {
			handler.OnError(err)
			return
		}
		switch UserDataEventType(gjson.GetBytes(*message, "e").String()) {
		case OutboundAccountPosition:
			var accountUpdateEvent WsAccountUpdateEvent
			err = s.c.GetHttpClient().JSONUnmarshal(*message, &accountUpdateEvent)
			if err != nil {
				handler.OnError(err)
				return
			}
			handler.OnAccountUpdate(&accountUpdateEvent)
		case ExecutionReport:
			var orderUpdateEvent WsOrderUpdateEvent
			err = s.c.GetHttpClient().JSONUnmarshal(*message, &orderUpdateEvent)
			if err != nil {
				handler.OnError(err)
				return
			}
			handler.OnOrderUpdate(&orderUpdateEvent)
		case ListenedKeyExpired:
			var listenKeyExpiredEvent WsListenKeyExpiredEvent
			err = s.c.GetHttpClient().JSONUnmarshal(*message, &listenKeyExpiredEvent)
			if err != nil {
				handler.OnError(err)
				return
			}
			handler.OnListenKeyExpired(&listenKeyExpiredEvent)
		default:
			handler.OnError(fmt.Errorf("unknown event type: %s", gjson.GetBytes(*message, "e").String()))
			return
		}
	}
	return request.Subscribe(ctx, s.c, s.listenKey, callback)
}

type UserDataHandler interface {
	OnAccountUpdate(message *WsAccountUpdateEvent)
	OnOrderUpdate(message *WsOrderUpdateEvent)
	OnListenKeyExpired(message *WsListenKeyExpiredEvent)
	OnError(err error)
}

type WsAccountUpdateEvent struct {
	Event             UserDataEventType  `json:"e"`
	Time              time.Time          `json:"E,format:unixmilli"`
	AccountUpdateTime time.Time          `json:"u,format:unix"`
	TransactionTime   time.Time          `json:"T,format:unix"`
	Balances          []WsAccountBalance `json:"B"`
	EventReasonType   string             `json:"m"`
}

type WsAccountBalance struct {
	Asset  string          `json:"a"`
	Free   decimal.Decimal `json:"f"`
	Locked decimal.Decimal `json:"l"`
}

type WsOrderUpdateEvent struct {
	Event                                  UserDataEventType `json:"e"`
	Time                                   time.Time         `json:"E,format:unixmilli"`
	Symbol                                 string            `json:"s"`
	ClientOrderId                          string            `json:"c"`
	OrderSide                              OrderSide         `json:"S"`
	OrderType                              OrderType         `json:"o"`
	TimeInForce                            TimeInForce       `json:"f"`
	Quantity                               decimal.Decimal   `json:"q"`
	Price                                  decimal.Decimal   `json:"p"`
	AveragePrice                           decimal.Decimal   `json:"ap"`
	StopPrice                              decimal.Decimal   `json:"P"`
	ExecutionType                          ExecutionType     `json:"x"`
	OrderStatus                            OrderStatus       `json:"X"`
	OrderId                                int64             `json:"i"`
	LastExecutedQuantity                   decimal.Decimal   `json:"l"`
	CumulativeFilledQuantity               decimal.Decimal   `json:"z"`
	LastExecutedPrice                      decimal.Decimal   `json:"L"`
	CommissionAmount                       decimal.Decimal   `json:"n"`
	CommissionAsset                        string            `json:"N"`
	TransactionTime                        time.Time         `json:"T,format:unixmilli"`
	TransactionId                          int64             `json:"t"`
	IsMaker                                bool              `json:"m"`
	OriginalOrderType                      OrderType         `json:"ot"`
	OrderCreationTime                      time.Time         `json:"O,format:unixmilli"`
	CumulativeQuoteAssetTransactedQuantity decimal.Decimal   `json:"Z"`
	LastQuoteAssetTransactedQuantity       decimal.Decimal   `json:"Y"`
	QuoteOrderQty                          decimal.Decimal   `json:"Q"`
}

type WsListenKeyExpiredEvent struct {
	Event UserDataEventType `json:"e"`
	Time  time.Time         `json:"E,format:unixmilli"`
}

type UserDataEventType string

const (
	OutboundAccountPosition UserDataEventType = "outboundAccountPosition"
	ExecutionReport         UserDataEventType = "executionReport"
	ListenedKeyExpired      UserDataEventType = "listenKeyExpired"
)

type ExecutionType string

const (
	NEW      ExecutionType = "NEW"
	CANCELED ExecutionType = "CANCELED"
	REJECTED ExecutionType = "REJECTED"
	TRADE    ExecutionType = "TRADE"
	EXPIRED  ExecutionType = "EXPIRED"
)
