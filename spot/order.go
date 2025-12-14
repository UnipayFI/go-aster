package spot

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type CreateOrderService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewCreateOrderService(symbol string, side OrderSide, orderType OrderType) *CreateOrderService {
	return &CreateOrderService{
		c:      c,
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

func (s *CreateOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/order", s.params).Sign()
	return request.Do[OrderResponse](req)
}

type OrderResponse struct {
	Symbol        string          `json:"symbol"`
	OrderId       int64           `json:"orderId"`
	ClientOrderId string          `json:"clientOrderId"`
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
	Time          time.Time       `json:"time,format:unixmilli"`
	UpdateTime    time.Time       `json:"updateTime,format:unixmilli"`
}

type CancelOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelOrderService(symbol string) *CancelOrderService {
	return &CancelOrderService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *CancelOrderService) SetOrderId(orderId int64) *CancelOrderService {
	s.params["orderId"] = strconv.FormatInt(orderId, 10)
	return s
}

func (s *CancelOrderService) SetOrigClientOrderId(origClientOrderId string) *CancelOrderService {
	s.params["origClientOrderId"] = origClientOrderId
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v1/order", s.params).Sign()
	return request.Do[OrderResponse](req)
}

type GetOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOrderService(symbol string) *GetOrderService {
	return &GetOrderService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *GetOrderService) SetOrderId(orderId int64) *GetOrderService {
	s.params["orderId"] = strconv.FormatInt(orderId, 10)
	return s
}

func (s *GetOrderService) SetOrigClientOrderId(origClientOrderId string) *GetOrderService {
	s.params["origClientOrderId"] = origClientOrderId
	return s
}

func (s *GetOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/order", s.params).Sign()
	return request.Do[OrderResponse](req)
}

type GetOpenOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{
		c:      c,
		params: map[string]string{},
	}
}

func (s *GetOpenOrdersService) SetSymbol(symbol string) *GetOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetOpenOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/openOrders", s.params).Sign()
	orders, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *orders, nil
}

type CancelAllOpenOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelAllOpenOrdersService(symbol string) *CancelAllOpenOrdersService {
	return &CancelAllOpenOrdersService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *CancelAllOpenOrdersService) SetOrderIdList(orderIds []int64) *CancelAllOpenOrdersService {
	data, _ := json.Marshal(orderIds)
	s.params["orderIdList"] = string(data)
	return s
}

func (s *CancelAllOpenOrdersService) SetOrigClientOrderIdList(origClientOrderIds []string) *CancelAllOpenOrdersService {
	data, _ := json.Marshal(origClientOrderIds)
	s.params["origClientOrderIdList"] = string(data)
	return s
}

func (s *CancelAllOpenOrdersService) Do(ctx context.Context) (*CancelAllOpenOrdersResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v1/allOpenOrders", s.params).Sign()
	return request.Do[CancelAllOpenOrdersResponse](req)
}

type CancelAllOpenOrdersResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetAllOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAllOrdersService(symbol string) *GetAllOrdersService {
	return &GetAllOrdersService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *GetAllOrdersService) SetOrderId(orderId int64) *GetAllOrdersService {
	s.params["orderId"] = strconv.FormatInt(orderId, 10)
	return s
}

func (s *GetAllOrdersService) SetStartTime(startTime time.Time) *GetAllOrdersService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetAllOrdersService) SetEndTime(endTime time.Time) *GetAllOrdersService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetAllOrdersService) SetLimit(limit int) *GetAllOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAllOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/allOrders", s.params).Sign()
	orders, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *orders, nil
}
