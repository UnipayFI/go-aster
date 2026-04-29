package spot

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// OrderResponse covers Place / Cancel / Get / OpenOrder responses; fields are
// a superset and may be empty for certain order types.
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

// PlaceOrderService -- POST /api/v3/order (TRADE)
type PlaceOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewPlaceOrderService(symbol string, side OrderSide, orderType OrderType) *PlaceOrderService {
	return &PlaceOrderService{c: c, params: map[string]string{
		"symbol": symbol,
		"side":   string(side),
		"type":   string(orderType),
	}}
}

func (s *PlaceOrderService) SetTimeInForce(tif TimeInForce) *PlaceOrderService {
	s.params["timeInForce"] = string(tif)
	return s
}

func (s *PlaceOrderService) SetQuantity(q decimal.Decimal) *PlaceOrderService {
	s.params["quantity"] = q.String()
	return s
}

func (s *PlaceOrderService) SetQuoteOrderQty(q decimal.Decimal) *PlaceOrderService {
	s.params["quoteOrderQty"] = q.String()
	return s
}

func (s *PlaceOrderService) SetPrice(p decimal.Decimal) *PlaceOrderService {
	s.params["price"] = p.String()
	return s
}

func (s *PlaceOrderService) SetNewClientOrderId(id string) *PlaceOrderService {
	s.params["newClientOrderId"] = id
	return s
}

func (s *PlaceOrderService) SetStopPrice(p decimal.Decimal) *PlaceOrderService {
	s.params["stopPrice"] = p.String()
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Post(s.c, ctx, "/api/v3/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

// CancelOrderService -- DELETE /api/v3/order (TRADE)
type CancelOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelOrderService(symbol string) *CancelOrderService {
	return &CancelOrderService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *CancelOrderService) SetOrderId(id int64) *CancelOrderService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *CancelOrderService) SetOrigClientOrderId(id string) *CancelOrderService {
	s.params["origClientOrderId"] = id
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v3/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

// GetOrderService -- GET /api/v3/order (USER_DATA)
type GetOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOrderService(symbol string) *GetOrderService {
	return &GetOrderService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetOrderService) SetOrderId(id int64) *GetOrderService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetOrderService) SetOrigClientOrderId(id string) *GetOrderService {
	s.params["origClientOrderId"] = id
	return s
}

func (s *GetOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/order", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

// GetOpenOrderService -- GET /api/v3/openOrder (USER_DATA)
//
// Returns a single open order; useful for fast lookup of one specific order.
// Use GetOpenOrdersService to list all.
type GetOpenOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOpenOrderService(symbol string) *GetOpenOrderService {
	return &GetOpenOrderService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetOpenOrderService) SetOrderId(id int64) *GetOpenOrderService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetOpenOrderService) SetOrigClientOrderId(id string) *GetOpenOrderService {
	s.params["origClientOrderId"] = id
	return s
}

func (s *GetOpenOrderService) Do(ctx context.Context) (*OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/openOrder", s.params).WithSignature()
	return request.Do[OrderResponse](req)
}

// GetOpenOrdersService -- GET /api/v3/openOrders (USER_DATA)
type GetOpenOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: c, params: map[string]string{}}
}

func (s *GetOpenOrdersService) SetSymbol(symbol string) *GetOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetOpenOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/openOrders", s.params).WithSignature()
	resp, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CancelAllOpenOrdersService -- DELETE /api/v3/allOpenOrders (TRADE)
type CancelAllOpenOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelAllOpenOrdersService(symbol string) *CancelAllOpenOrdersService {
	return &CancelAllOpenOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *CancelAllOpenOrdersService) SetOrderIdList(ids []int64) *CancelAllOpenOrdersService {
	data, _ := json.Marshal(ids)
	s.params["orderIdList"] = string(data)
	return s
}

func (s *CancelAllOpenOrdersService) SetOrigClientOrderIdList(ids []string) *CancelAllOpenOrdersService {
	data, _ := json.Marshal(ids)
	s.params["origClientOrderIdList"] = string(data)
	return s
}

func (s *CancelAllOpenOrdersService) Do(ctx context.Context) (*CancelAllOpenOrdersResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v3/allOpenOrders", s.params).WithSignature()
	return request.Do[CancelAllOpenOrdersResponse](req)
}

type CancelAllOpenOrdersResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// GetAllOrdersService -- GET /api/v3/allOrders (USER_DATA)
type GetAllOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAllOrdersService(symbol string) *GetAllOrdersService {
	return &GetAllOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetAllOrdersService) SetOrderId(id int64) *GetAllOrdersService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetAllOrdersService) SetStartTime(t time.Time) *GetAllOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAllOrdersService) SetEndTime(t time.Time) *GetAllOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAllOrdersService) SetLimit(limit int) *GetAllOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAllOrdersService) Do(ctx context.Context) ([]OrderResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/allOrders", s.params).WithSignature()
	resp, err := request.Do[[]OrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetTransactionHistoryService -- GET /api/v3/transactionHistory (USER_DATA)
type GetTransactionHistoryService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTransactionHistoryService() *GetTransactionHistoryService {
	return &GetTransactionHistoryService{c: c, params: map[string]string{}}
}

func (s *GetTransactionHistoryService) SetAsset(asset string) *GetTransactionHistoryService {
	s.params["asset"] = asset
	return s
}

func (s *GetTransactionHistoryService) SetType(t TransactionType) *GetTransactionHistoryService {
	s.params["type"] = string(t)
	return s
}

func (s *GetTransactionHistoryService) SetStartTime(t time.Time) *GetTransactionHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetTransactionHistoryService) SetEndTime(t time.Time) *GetTransactionHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetTransactionHistoryService) SetLimit(limit int) *GetTransactionHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetTransactionHistoryService) Do(ctx context.Context) ([]TransactionRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/transactionHistory", s.params).WithSignature()
	resp, err := request.Do[[]TransactionRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type TransactionRecord struct {
	TranId        int64           `json:"tranId"`
	TradeId       *int64          `json:"tradeId"`
	Asset         string          `json:"asset"`
	Symbol        string          `json:"symbol"`
	BalanceDelta  decimal.Decimal `json:"balanceDelta"`
	BalanceInfo   string          `json:"balanceInfo"`
	Time          time.Time       `json:"time,format:unixmilli"`
	Type          TransactionType `json:"type"`
}
