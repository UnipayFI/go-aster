package futures

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// Order is the response shape for new/modify/cancel/query order endpoints.
// TRAILING_STOP_MARKET-only fields (ActivatePrice, PriceRate) are populated
// only for that order type; other fields are common.
type Order struct {
	OrderId       int64           `json:"orderId"`
	Symbol        string          `json:"symbol"`
	Status        OrderStatus     `json:"status"`
	ClientOrderId string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	OrigQty       decimal.Decimal `json:"origQty"`
	ExecutedQty   decimal.Decimal `json:"executedQty"`
	CumQty        decimal.Decimal `json:"cumQty"`
	CumQuote      decimal.Decimal `json:"cumQuote"`
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

// PlaceOrderService -- POST /fapi/v3/order (TRADE)
type PlaceOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewPlaceOrderService(symbol string, side OrderSide, orderType OrderType) *PlaceOrderService {
	return &PlaceOrderService{c: c, params: map[string]string{
		"symbol": symbol,
		"side":   string(side),
		"type":   string(orderType),
	}}
}

func (s *PlaceOrderService) SetPositionSide(p PositionSide) *PlaceOrderService {
	s.params["positionSide"] = string(p)
	return s
}

func (s *PlaceOrderService) SetTimeInForce(tif TimeInForce) *PlaceOrderService {
	s.params["timeInForce"] = string(tif)
	return s
}

func (s *PlaceOrderService) SetQuantity(q decimal.Decimal) *PlaceOrderService {
	s.params["quantity"] = q.String()
	return s
}

func (s *PlaceOrderService) SetReduceOnly(reduceOnly bool) *PlaceOrderService {
	s.params["reduceOnly"] = strconv.FormatBool(reduceOnly)
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

func (s *PlaceOrderService) SetClosePosition(close bool) *PlaceOrderService {
	s.params["closePosition"] = strconv.FormatBool(close)
	return s
}

func (s *PlaceOrderService) SetActivationPrice(p decimal.Decimal) *PlaceOrderService {
	s.params["activationPrice"] = p.String()
	return s
}

func (s *PlaceOrderService) SetCallbackRate(r decimal.Decimal) *PlaceOrderService {
	s.params["callbackRate"] = r.String()
	return s
}

func (s *PlaceOrderService) SetWorkingType(w WorkingType) *PlaceOrderService {
	s.params["workingType"] = string(w)
	return s
}

func (s *PlaceOrderService) SetPriceProtect(p bool) *PlaceOrderService {
	if p {
		s.params["priceProtect"] = "TRUE"
	} else {
		s.params["priceProtect"] = "FALSE"
	}
	return s
}

func (s *PlaceOrderService) SetNewOrderRespType(r ResponseType) *PlaceOrderService {
	s.params["newOrderRespType"] = string(r)
	return s
}

// SetPegPriceType turns a LIMIT order into a BBO-pegged order; the engine
// resolves the actual price from the order book at trigger time using the
// chosen BBO level plus the pegOffset.
func (s *PlaceOrderService) SetPegPriceType(p PegPriceType) *PlaceOrderService {
	s.params["pegPriceType"] = string(p)
	return s
}

// SetPegOffset sets the signed offset from the BBO. BUY orders use a
// non-positive value, SELL a non-negative value; must be a tickSize multiple.
func (s *PlaceOrderService) SetPegOffset(o decimal.Decimal) *PlaceOrderService {
	s.params["pegOffset"] = o.String()
	return s
}

// SetSTPMode overrides the account-level Self-Trade Prevention mode for this
// single order.
func (s *PlaceOrderService) SetSTPMode(m STPMode) *PlaceOrderService {
	s.params["stpMode"] = string(m)
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/order", s.params).WithSignature()
	return request.Do[Order](req)
}

// ModifyOrderService -- PUT /fapi/v3/order (TRADE)
//
// Only LIMIT orders can be modified, and only quantity / price. Either
// orderId or origClientOrderId must be set; orderId wins if both.
type ModifyOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewModifyOrderService(symbol string) *ModifyOrderService {
	return &ModifyOrderService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *ModifyOrderService) SetOrderId(id int64) *ModifyOrderService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *ModifyOrderService) SetOrigClientOrderId(id string) *ModifyOrderService {
	s.params["origClientOrderId"] = id
	return s
}

func (s *ModifyOrderService) SetQuantity(q decimal.Decimal) *ModifyOrderService {
	s.params["quantity"] = q.String()
	return s
}

func (s *ModifyOrderService) SetPrice(p decimal.Decimal) *ModifyOrderService {
	s.params["price"] = p.String()
	return s
}

func (s *ModifyOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Put(ctx, s.c, "/fapi/v3/order", s.params).WithSignature()
	return request.Do[Order](req)
}

// PlaceChaseOrderService -- POST /fapi/v3/chase (TRADE)
//
// A Chase order is a BBO-pegged GTX (post-only) limit order. The strategy
// service re-pegs the order to bid1-chaseOffset (BUY) or ask1+chaseOffset
// (SELL) each tick until it fills, is cancelled (via DELETE /fapi/v3/order),
// or the market moves beyond maxChaseOffset from the original BBO. To stop a
// BBO-pegged order from re-resolving you must use this endpoint rather than a
// plain modify.
//
// Note: chaseOffsetType only supports ABSOLUTE for now; PERCENTAGE is rejected
// with UNSUPPORTED_OPERATION. maxChaseOffsetType supports both.
type PlaceChaseOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewPlaceChaseOrderService(symbol string, side OrderSide, quantityUnit QuantityUnit, quantity decimal.Decimal) *PlaceChaseOrderService {
	return &PlaceChaseOrderService{c: c, params: map[string]string{
		"symbol":       symbol,
		"side":         string(side),
		"quantityUnit": string(quantityUnit),
		"quantity":     quantity.String(),
	}}
}

func (s *PlaceChaseOrderService) SetPositionSide(p PositionSide) *PlaceChaseOrderService {
	s.params["positionSide"] = string(p)
	return s
}

func (s *PlaceChaseOrderService) SetReduceOnly(reduceOnly bool) *PlaceChaseOrderService {
	s.params["reduceOnly"] = strconv.FormatBool(reduceOnly)
	return s
}

func (s *PlaceChaseOrderService) SetChaseOffset(o decimal.Decimal) *PlaceChaseOrderService {
	s.params["chaseOffset"] = o.String()
	return s
}

func (s *PlaceChaseOrderService) SetChaseOffsetType(t OffsetType) *PlaceChaseOrderService {
	s.params["chaseOffsetType"] = string(t)
	return s
}

func (s *PlaceChaseOrderService) SetMaxChaseOffset(o decimal.Decimal) *PlaceChaseOrderService {
	s.params["maxChaseOffset"] = o.String()
	return s
}

func (s *PlaceChaseOrderService) SetMaxChaseOffsetType(t OffsetType) *PlaceChaseOrderService {
	s.params["maxChaseOffsetType"] = string(t)
	return s
}

func (s *PlaceChaseOrderService) SetTimeInForce(tif TimeInForce) *PlaceChaseOrderService {
	s.params["timeInForce"] = string(tif)
	return s
}

func (s *PlaceChaseOrderService) SetClientStrategyId(id string) *PlaceChaseOrderService {
	s.params["clientStrategyId"] = id
	return s
}

func (s *PlaceChaseOrderService) Do(ctx context.Context) (*ChaseOrder, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/chase", s.params).WithSignature()
	return request.Do[ChaseOrder](req)
}

// ChaseOrder is the response shape for a placed chase strategy order.
type ChaseOrder struct {
	StrategyId         int64           `json:"strategyId"`
	ClientStrategyId   string          `json:"clientStrategyId"`
	Symbol             string          `json:"symbol"`
	Side               OrderSide       `json:"side"`
	PositionSide       PositionSide    `json:"positionSide"`
	Quantity           decimal.Decimal `json:"quantity"`
	QuantityUnit       QuantityUnit    `json:"quantityUnit"`
	ReduceOnly         bool            `json:"reduceOnly"`
	ChaseOffset        decimal.Decimal `json:"chaseOffset"`
	ChaseOffsetType    OffsetType      `json:"chaseOffsetType"`
	MaxChaseOffset     decimal.Decimal `json:"maxChaseOffset"`
	MaxChaseOffsetType OffsetType      `json:"maxChaseOffsetType"`
	TimeInForce        TimeInForce     `json:"timeInForce"`
	StrategyStatus     string          `json:"strategyStatus"`
	BookTime           time.Time       `json:"bookTime,format:unixmilli"`
	UpdateTime         time.Time       `json:"updateTime,format:unixmilli"`
}

// BatchOrderItem mirrors the per-order parameters accepted in batchOrders.
// Empty string fields are omitted from the marshalled JSON.
type BatchOrderItem struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	PositionSide     string `json:"positionSide,omitempty"`
	Type             string `json:"type"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	Quantity         string `json:"quantity"`
	ReduceOnly       string `json:"reduceOnly,omitempty"`
	Price            string `json:"price,omitempty"`
	NewClientOrderId string `json:"newClientOrderId,omitempty"`
	StopPrice        string `json:"stopPrice,omitempty"`
	ActivationPrice  string `json:"activationPrice,omitempty"`
	CallbackRate     string `json:"callbackRate,omitempty"`
	WorkingType      string `json:"workingType,omitempty"`
	PriceProtect     string `json:"priceProtect,omitempty"`
	NewOrderRespType string `json:"newOrderRespType,omitempty"`
}

// BatchOrdersService -- POST /fapi/v3/batchOrders (TRADE)
//
// Each list entry returns either an Order or {code,msg} error. Inspect the
// raw responses with json.RawMessage if you need to handle partial failures.
type BatchOrdersService struct {
	c      *FuturesClient
	orders []BatchOrderItem
}

func (c *FuturesClient) NewBatchOrdersService() *BatchOrdersService {
	return &BatchOrdersService{c: c}
}

func (s *BatchOrdersService) SetOrders(orders []BatchOrderItem) *BatchOrdersService {
	s.orders = orders
	return s
}

func (s *BatchOrdersService) Do(ctx context.Context) ([]json.RawMessage, error) {
	data, err := json.Marshal(s.orders)
	if err != nil {
		return nil, err
	}
	req := request.Post(s.c, ctx, "/fapi/v3/batchOrders", map[string]string{
		"batchOrders": string(data),
	}).WithSignature()
	resp, err := request.Do[[]json.RawMessage](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchModifyOrderItem mirrors the per-order parameters accepted in a batch
// modify. Empty string fields are omitted from the marshalled JSON. Either
// OrderId or OrigClientOrderId must be set; Quantity and Price are both
// required even when only one of them changes.
type BatchModifyOrderItem struct {
	Symbol            string `json:"symbol"`
	OrderId           string `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
	Quantity          string `json:"quantity"`
	Price             string `json:"price"`
}

// ModifyMultipleOrdersService -- PUT /fapi/v3/batchOrders (TRADE)
//
// Modifies up to 5 orders in one call (10 for market-maker whitelisted
// accounts). Per-order rules match ModifyOrderService: LIMIT orders only, and
// at most 10000 modifications per individual order.
//
// Orders are processed concurrently and independently -- there is no atomic
// guarantee, so one entry failing leaves the others untouched -- but the
// response list keeps the request order. Each entry is either an Order or
// {code,msg}; inspect the raw messages to handle partial failures.
type ModifyMultipleOrdersService struct {
	c      *FuturesClient
	orders []BatchModifyOrderItem
}

func (c *FuturesClient) NewModifyMultipleOrdersService() *ModifyMultipleOrdersService {
	return &ModifyMultipleOrdersService{c: c}
}

func (s *ModifyMultipleOrdersService) SetOrders(orders []BatchModifyOrderItem) *ModifyMultipleOrdersService {
	s.orders = orders
	return s
}

func (s *ModifyMultipleOrdersService) Do(ctx context.Context) ([]json.RawMessage, error) {
	data, err := json.Marshal(s.orders)
	if err != nil {
		return nil, err
	}
	req := request.Put(ctx, s.c, "/fapi/v3/batchOrders", map[string]string{
		"batchOrders": string(data),
	}).WithSignature()
	resp, err := request.Do[[]json.RawMessage](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesSpotTransferService -- POST /fapi/v3/asset/wallet/transfer (TRANSFER)
type FuturesSpotTransferService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewFuturesSpotTransferService(asset string, amount decimal.Decimal, kind TransferKindType, clientTranID string) *FuturesSpotTransferService {
	return &FuturesSpotTransferService{c: c, params: map[string]string{
		"amount":       amount.String(),
		"asset":        asset,
		"clientTranId": clientTranID,
		"kindType":     string(kind),
	}}
}

func (s *FuturesSpotTransferService) Do(ctx context.Context) (*TransferResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/asset/wallet/transfer", s.params).WithSignature()
	return request.Do[TransferResponse](req)
}

type TransferResponse struct {
	TranID int64  `json:"tranId"`
	Status string `json:"status"`
}

// GetOrderService -- GET /fapi/v3/order (USER_DATA)
type GetOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOrderService(symbol string) *GetOrderService {
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

func (s *GetOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/order", s.params).WithSignature()
	return request.Do[Order](req)
}

// CancelOrderService -- DELETE /fapi/v3/order (TRADE)
type CancelOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCancelOrderService(symbol string) *CancelOrderService {
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

func (s *CancelOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/order", s.params).WithSignature()
	return request.Do[Order](req)
}

// CancelAllOpenOrdersService -- DELETE /fapi/v3/allOpenOrders (TRADE)
type CancelAllOpenOrdersService struct {
	c      *FuturesClient
	symbol string
}

func (c *FuturesClient) NewCancelAllOpenOrdersService(symbol string) *CancelAllOpenOrdersService {
	return &CancelAllOpenOrdersService{c: c, symbol: symbol}
}

func (s *CancelAllOpenOrdersService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/allOpenOrders", map[string]string{"symbol": s.symbol}).WithSignature()
	return request.Do[GenericCodeMsg](req)
}

type GenericCodeMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CancelMultipleOrdersService -- DELETE /fapi/v3/batchOrders (TRADE)
type CancelMultipleOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewCancelMultipleOrdersService(symbol string) *CancelMultipleOrdersService {
	return &CancelMultipleOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *CancelMultipleOrdersService) SetOrderIdList(ids []int64) *CancelMultipleOrdersService {
	data, _ := json.Marshal(ids)
	s.params["orderIdList"] = string(data)
	return s
}

func (s *CancelMultipleOrdersService) SetOrigClientOrderIdList(ids []string) *CancelMultipleOrdersService {
	data, _ := json.Marshal(ids)
	s.params["origClientOrderIdList"] = string(data)
	return s
}

func (s *CancelMultipleOrdersService) Do(ctx context.Context) ([]json.RawMessage, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/batchOrders", s.params).WithSignature()
	resp, err := request.Do[[]json.RawMessage](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GuardedCancelOrderService -- DELETE /fapi/v3/guardedCancelOrder (TRADE)
//
// Same parameters and response as DELETE /fapi/v3/order, except the signature
// must replay the exact nonce that was submitted when the order was placed
// (POST /fapi/v3/order), not a fresh timestamp. The engine uses that on-chain
// nonce to guard against duplicate or out-of-order cancellations. Either
// orderId or origClientOrderId must be set.
type GuardedCancelOrderService struct {
	c      *FuturesClient
	nonce  int64
	params map[string]string
}

func (c *FuturesClient) NewGuardedCancelOrderService(symbol string, nonce int64) *GuardedCancelOrderService {
	return &GuardedCancelOrderService{c: c, nonce: nonce, params: map[string]string{"symbol": symbol}}
}

func (s *GuardedCancelOrderService) SetOrderId(id int64) *GuardedCancelOrderService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GuardedCancelOrderService) SetOrigClientOrderId(id string) *GuardedCancelOrderService {
	s.params["origClientOrderId"] = id
	return s
}

func (s *GuardedCancelOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/guardedCancelOrder", s.params).WithSignatureNonce(s.nonce)
	return request.Do[Order](req)
}

// GuardedBatchOrdersService -- DELETE /fapi/v3/guardedBatchOrders (TRADE)
//
// Same parameters and response as DELETE /fapi/v3/batchOrders (cancel up to 10
// orders), except the signature must replay the exact nonce submitted when the
// orders were placed (POST /fapi/v3/batchOrders), not a fresh timestamp.
type GuardedBatchOrdersService struct {
	c      *FuturesClient
	nonce  int64
	params map[string]string
}

func (c *FuturesClient) NewGuardedBatchOrdersService(symbol string, nonce int64) *GuardedBatchOrdersService {
	return &GuardedBatchOrdersService{c: c, nonce: nonce, params: map[string]string{"symbol": symbol}}
}

func (s *GuardedBatchOrdersService) SetOrderIdList(ids []int64) *GuardedBatchOrdersService {
	data, _ := json.Marshal(ids)
	s.params["orderIdList"] = string(data)
	return s
}

func (s *GuardedBatchOrdersService) SetOrigClientOrderIdList(ids []string) *GuardedBatchOrdersService {
	data, _ := json.Marshal(ids)
	s.params["origClientOrderIdList"] = string(data)
	return s
}

func (s *GuardedBatchOrdersService) Do(ctx context.Context) ([]json.RawMessage, error) {
	req := request.Delete(ctx, s.c, "/fapi/v3/guardedBatchOrders", s.params).WithSignatureNonce(s.nonce)
	resp, err := request.Do[[]json.RawMessage](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CountdownCancelAllService -- POST /fapi/v3/countdownCancelAll (TRADE)
//
// Heartbeat-style: call repeatedly to keep extending the countdown. A value
// of 0 disables the timer. The system polls roughly every 10ms, so don't
// rely on sub-second precision.
type CountdownCancelAllService struct {
	c             *FuturesClient
	symbol        string
	countdownTime int64
}

func (c *FuturesClient) NewCountdownCancelAllService(symbol string, countdownTimeMs int64) *CountdownCancelAllService {
	return &CountdownCancelAllService{c: c, symbol: symbol, countdownTime: countdownTimeMs}
}

func (s *CountdownCancelAllService) Do(ctx context.Context) (*CountdownCancelAllResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/countdownCancelAll", map[string]string{
		"symbol":        s.symbol,
		"countdownTime": strconv.FormatInt(s.countdownTime, 10),
	}).WithSignature()
	return request.Do[CountdownCancelAllResponse](req)
}

type CountdownCancelAllResponse struct {
	Symbol        string `json:"symbol"`
	CountdownTime string `json:"countdownTime"`
}

// GetOpenOrderService -- GET /fapi/v3/openOrder (USER_DATA)
type GetOpenOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOpenOrderService(symbol string) *GetOpenOrderService {
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

func (s *GetOpenOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/openOrder", s.params).WithSignature()
	return request.Do[Order](req)
}

// GetOpenOrdersService -- GET /fapi/v3/openOrders (USER_DATA)
type GetOpenOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: c, params: map[string]string{}}
}

func (s *GetOpenOrdersService) SetSymbol(symbol string) *GetOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetOpenOrdersService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/openOrders", s.params).WithSignature()
	resp, err := request.Do[[]Order](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetAllOrdersService -- GET /fapi/v3/allOrders (USER_DATA)
type GetAllOrdersService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetAllOrdersService(symbol string) *GetAllOrdersService {
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

func (s *GetAllOrdersService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/allOrders", s.params).WithSignature()
	resp, err := request.Do[[]Order](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
