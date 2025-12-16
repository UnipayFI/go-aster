package futures

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type CreateOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCreateOrderService(symbol string, side OrderSide, orderType OrderType) *CreateOrderService {
	return &CreateOrderService{
		c:      c,
		params: map[string]string{"symbol": symbol, "side": string(side), "type": string(orderType)},
	}
}

func (s *CreateOrderService) SetPositionSide(positionSide PositionSide) *CreateOrderService {
	s.params["positionSide"] = string(positionSide)
	return s
}

func (s *CreateOrderService) SetTimeInForce(timeInForce TimeInForce) *CreateOrderService {
	s.params["timeInForce"] = string(timeInForce)
	return s
}

func (s *CreateOrderService) SetQuantity(quantity float64) *CreateOrderService {
	s.params["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetReduceOnly(reduceOnly bool) *CreateOrderService {
	s.params["reduceOnly"] = strconv.FormatBool(reduceOnly)
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

func (s *CreateOrderService) SetClosePosition(closePosition bool) *CreateOrderService {
	s.params["closePosition"] = strconv.FormatBool(closePosition)
	return s
}

func (s *CreateOrderService) SetActivationPrice(activationPrice float64) *CreateOrderService {
	s.params["activationPrice"] = strconv.FormatFloat(activationPrice, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetCallbackRate(callbackRate float64) *CreateOrderService {
	s.params["callbackRate"] = strconv.FormatFloat(callbackRate, 'f', -1, 64)
	return s
}

func (s *CreateOrderService) SetWorkingType(workingType WorkingType) *CreateOrderService {
	s.params["workingType"] = string(workingType)
	return s
}

func (s *CreateOrderService) SetPriceProtect(priceProtect bool) *CreateOrderService {
	s.params["priceProtect"] = strconv.FormatBool(priceProtect)
	return s
}

func (s *CreateOrderService) SetNewOrderRespType(newOrderRespType NewOrderRespType) *CreateOrderService {
	s.params["newOrderRespType"] = string(newOrderRespType)
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

type OrderResponse struct {
	Symbol        string          `json:"symbol"`
	OrderId       int64           `json:"orderId"`
	ClientOrderId string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	OrigQty       decimal.Decimal `json:"origQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	CumQty        decimal.Decimal `json:"cumQty"`
	CumQuote      decimal.Decimal `json:"cumQuote"`
	Status        OrderStatus     `json:"status"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	Type          OrderType       `json:"type"`
	OrigType      OrderType       `json:"origType"`
	ReduceOnly    bool            `json:"reduceOnly"`
	ClosePosition bool            `json:"closePosition"`
	Side          OrderSide       `json:"side"`
	PositionSide  PositionSide    `json:"positionSide"`
	StopPrice     decimal.Decimal `json:"stopPrice"`
	WorkingType   WorkingType     `json:"workingType"`
	PriceProtect  bool            `json:"priceProtect"`
	ActivatePrice decimal.Decimal `json:"activatePrice"`
	PriceRate     decimal.Decimal `json:"priceRate"`
	Time          time.Time       `json:"time,format:unixmilli"`
	UpdateTime    time.Time       `json:"updateTime,format:unixmilli"`
}

type CreateBatchOrdersService struct {
	c      *FuturesClient
	orders []map[string]string
}

func (c *FuturesClient) NewCreateBatchOrdersService() *CreateBatchOrdersService {
	return &CreateBatchOrdersService{c: c, orders: []map[string]string{}}
}

func (s *CreateBatchOrdersService) AddOrder(order map[string]string) *CreateBatchOrdersService {
	s.orders = append(s.orders, order)
	return s
}

func (s *CreateBatchOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	data, _ := json.Marshal(s.orders)
	params := map[string]string{"batchOrders": string(data)}
	req := request.Post(s.c, ctx, "/fapi/v1/batchOrders", params).WithSignature()
	resp, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CancelOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCancelOrderService(symbol string) *CancelOrderService {
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
	req := request.Delete(ctx, s.c, "/fapi/v1/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

type CancelAllOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCancelAllOrdersService(symbol string) *CancelAllOrdersService {
	return &CancelAllOrdersService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *CancelAllOrdersService) Do(ctx context.Context) (*CancelAllOrdersResponse, error) {
	req := request.Delete(ctx, s.c, "/fapi/v1/allOpenOrders", s.params).WithSignature()
	return request.Do[CancelAllOrdersResponse](req)
}

type CancelAllOrdersResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CancelBatchOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCancelBatchOrdersService(symbol string) *CancelBatchOrdersService {
	return &CancelBatchOrdersService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *CancelBatchOrdersService) SetOrderIdList(orderIds []int64) *CancelBatchOrdersService {
	data, _ := json.Marshal(orderIds)
	s.params["orderIdList"] = string(data)
	return s
}

func (s *CancelBatchOrdersService) SetOrigClientOrderIdList(origClientOrderIds []string) *CancelBatchOrdersService {
	data, _ := json.Marshal(origClientOrderIds)
	s.params["origClientOrderIdList"] = string(data)
	return s
}

func (s *CancelBatchOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	req := request.Delete(ctx, s.c, "/fapi/v1/batchOrders", s.params).WithSignature()
	resp, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type GetOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOrderService(symbol string) *GetOrderService {
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
	req := request.Get(ctx, s.c, "/fapi/v1/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

type GetOpenOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOpenOrderService(symbol string) *GetOpenOrderService {
	return &GetOpenOrderService{
		c:      c,
		params: map[string]string{"symbol": symbol},
	}
}

func (s *GetOpenOrderService) SetOrderId(orderId int64) *GetOpenOrderService {
	s.params["orderId"] = strconv.FormatInt(orderId, 10)
	return s
}

func (s *GetOpenOrderService) SetOrigClientOrderId(origClientOrderId string) *GetOpenOrderService {
	s.params["origClientOrderId"] = origClientOrderId
	return s
}

func (s *GetOpenOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Get(ctx, s.c, "/fapi/v1/openOrder", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

type GetOpenOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOpenOrdersService() *GetOpenOrdersService {
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
	req := request.Get(ctx, s.c, "/fapi/v1/openOrders", s.params).WithSignature()
	orders, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *orders, nil
}

type GetAllOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetAllOrdersService(symbol string) *GetAllOrdersService {
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
	req := request.Get(ctx, s.c, "/fapi/v1/allOrders", s.params).WithSignature()
	orders, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *orders, nil
}

type CountdownCancelAllService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCountdownCancelAllService(symbol string, countdownTime int64) *CountdownCancelAllService {
	return &CountdownCancelAllService{
		c: c,
		params: map[string]string{
			"symbol":        symbol,
			"countdownTime": strconv.FormatInt(countdownTime, 10),
		},
	}
}

func (s *CountdownCancelAllService) Do(ctx context.Context) (*CountdownCancelAllResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/countdownCancelAll", s.params).WithSignature()
	return request.Do[CountdownCancelAllResponse](req)
}

type CountdownCancelAllResponse struct {
	Symbol        string `json:"symbol"`
	CountdownTime string `json:"countdownTime"`
}

type WalletTransferService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewWalletTransferService(asset string, amount float64, clientTranId string, kindType TransferType) *WalletTransferService {
	return &WalletTransferService{
		c: c,
		params: map[string]string{
			"asset":        asset,
			"amount":       strconv.FormatFloat(amount, 'f', -1, 64),
			"clientTranId": clientTranId,
			"kindType":     string(kindType),
		},
	}
}

func (s *WalletTransferService) Do(ctx context.Context) (*WalletTransferResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v1/asset/wallet/transfer", s.params).WithSignature()
	return request.Do[WalletTransferResponse](req)
}

type WalletTransferResponse struct {
	TranId int64  `json:"tranId"`
	Status string `json:"status"`
}
