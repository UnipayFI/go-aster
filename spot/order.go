package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type CreateOrderService struct {
	client *SpotClient

	params map[string]string
}

func NewCreateOrderService(client *SpotClient, symbol string, side OrderSide, orderType OrderType) *CreateOrderService {
	return &CreateOrderService{
		client: client,
		params: map[string]string{"symbol": symbol, "side": string(side), "type": string(orderType)},
	}
}

func (s *CreateOrderService) SetTimeInForce(timeInForce TimeInForce) *CreateOrderService {
	s.params["timeInForce"] = string(timeInForce)
	return s
}

func (s *CreateOrderService) SetQuantity(quantity float64) *CreateOrderService {
	s.params["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetQuoteOrderQty(quoteOrderQty float64) *CreateOrderService {
	s.params["quoteOrderQty"] = strconv.FormatFloat(quoteOrderQty, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetPrice(price float64) *CreateOrderService {
	s.params["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetNewClientOrderId(newClientOrderId string) *CreateOrderService {
	s.params["newClientOrderId"] = newClientOrderId
	return s
}

func (s *CreateOrderService) SetStopPrice(stopPrice float64) *CreateOrderService {
	s.params["stopPrice"] = strconv.FormatFloat(stopPrice, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*CreateOrderResponse, error) {
	req := request.Post(s.client, ctx, "/api/v3/order", s.params).Sign()
	return request.Do[CreateOrderResponse](req)
}

type CreateOrderResponse struct {
	Symbol        string          `json:"symbol"`
	OrderId       int64           `json:"orderId"`
	ClientOrderId string          `json:"clientOrderId"`
	UpdateTime    time.Time       `json:"updateTime,format:unixmilli"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	OrigQty       decimal.Decimal `json:"origQty"`
	CumQty        decimal.Decimal `json:"cumQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	CumQuote      decimal.Decimal `json:"cumQuote"`
	Status        OrderStatus     `json:"status"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	StopPrice     decimal.Decimal `json:"stopPrice"`
	OrigType      OrderType       `json:"origType"`
	Type          OrderType       `json:"type"`
	Side          OrderSide       `json:"side"`
}
