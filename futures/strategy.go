package futures

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// StrategySubOrder is one leg of a strategy order. The same shape is used for
// both placeStrategyOrder and updateStrategyOrder; empty-string fields are
// omitted from the marshalled JSON so the engine applies its defaults.
//
// strategySubId is the 1-based sequence number and must match the element's
// position in the list. The firstDriven*/secondDriven*/`*Trigger` fields wire
// up the conditional logic (e.g. for OTO the second order is PLACE_ORDER-ed
// when the first one is FILLED).
type StrategySubOrder struct {
	StrategySubId   string `json:"strategySubId"`
	SecurityType    string `json:"securityType"`
	Symbol          string `json:"symbol"`
	Side            string `json:"side"`
	PositionSide    string `json:"positionSide,omitempty"`
	Type            string `json:"type"`
	Quantity        string `json:"quantity,omitempty"`
	Price           string `json:"price,omitempty"`
	StopPrice       string `json:"stopPrice,omitempty"`
	TimeInForce     string `json:"timeInForce,omitempty"`
	WorkingType     string `json:"workingType,omitempty"`
	ReduceOnly      string `json:"reduceOnly,omitempty"`
	ClosePosition   string `json:"closePosition,omitempty"`
	PriceProtect    string `json:"priceProtect,omitempty"`
	ClientOrderId   string `json:"clientOrderId,omitempty"`
	ActivationPrice string `json:"activationPrice,omitempty"`
	CallbackRate    string `json:"callbackRate,omitempty"`
	FirstDrivenId   string `json:"firstDrivenId,omitempty"`
	FirstDrivenOn   string `json:"firstDrivenOn,omitempty"`
	FirstTrigger    string `json:"firstTrigger,omitempty"`
	SecondDrivenId  string `json:"secondDrivenId,omitempty"`
	SecondDrivenOn  string `json:"secondDrivenOn,omitempty"`
	SecondTrigger   string `json:"secondTrigger,omitempty"`
}

// PlaceStrategyOrderService -- POST /fapi/v3/placeStrategyOrder (TRADE)
//
// OTO and OCO require exactly 2 sub-orders; OTOCO requires exactly 3.
type PlaceStrategyOrderService struct {
	c                *FuturesClient
	strategyType     StrategyType
	clientStrategyId string
	subOrders        []StrategySubOrder
}

func (c *FuturesClient) NewPlaceStrategyOrderService(strategyType StrategyType, subOrders []StrategySubOrder) *PlaceStrategyOrderService {
	return &PlaceStrategyOrderService{c: c, strategyType: strategyType, subOrders: subOrders}
}

func (s *PlaceStrategyOrderService) SetClientStrategyId(id string) *PlaceStrategyOrderService {
	s.clientStrategyId = id
	return s
}

func (s *PlaceStrategyOrderService) Do(ctx context.Context) (*StrategyOrderResult, error) {
	data, err := json.Marshal(s.subOrders)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"strategyType": string(s.strategyType),
		"subOrderList": string(data),
	}
	if s.clientStrategyId != "" {
		params["clientStrategyId"] = s.clientStrategyId
	}
	req := request.Post(s.c, ctx, "/fapi/v3/placeStrategyOrder", params).WithSignature()
	return request.Do[StrategyOrderResult](req)
}

// StrategyOrderResult is the placeStrategyOrder response.
type StrategyOrderResult struct {
	StrategyId       int64        `json:"strategyId"`
	ClientStrategyId string       `json:"clientStrategyId"`
	StrategyType     StrategyType `json:"strategyType"`
	StrategyStatus   string       `json:"strategyStatus"`
	UpdateTime       time.Time    `json:"updateTime,format:unixmilli"`
	FailureCode      int          `json:"failureCode"`
	FailureReason    string       `json:"failureReason"`
}

// UpdateStrategyOrderService -- POST /fapi/v3/updateStrategyOrder (TRADE)
//
// Updates one or more sub-orders of an existing strategy (max 2 for OTO/OCO,
// max 3 for OTOCO). Returns one result per updated sub-order.
type UpdateStrategyOrderService struct {
	c            *FuturesClient
	strategyId   int64
	strategyType StrategyType
	subOrders    []StrategySubOrder
}

func (c *FuturesClient) NewUpdateStrategyOrderService(strategyId int64, strategyType StrategyType, subOrders []StrategySubOrder) *UpdateStrategyOrderService {
	return &UpdateStrategyOrderService{c: c, strategyId: strategyId, strategyType: strategyType, subOrders: subOrders}
}

func (s *UpdateStrategyOrderService) Do(ctx context.Context) ([]StrategyUpdateResult, error) {
	data, err := json.Marshal(s.subOrders)
	if err != nil {
		return nil, err
	}
	req := request.Post(s.c, ctx, "/fapi/v3/updateStrategyOrder", map[string]string{
		"strategyId":   strconv.FormatInt(s.strategyId, 10),
		"strategyType": string(s.strategyType),
		"subOrderList": string(data),
	}).WithSignature()
	resp, err := request.Do[[]StrategyUpdateResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type StrategyUpdateResult struct {
	StrategyId       int64        `json:"strategyId"`
	ClientStrategyId string       `json:"clientStrategyId"`
	StrategyType     StrategyType `json:"strategyType"`
	StrategyStatus   string       `json:"strategyStatus"`
	UpdatedSubOrder  int          `json:"updatedSubOrder"`
	UpdateStatus     string       `json:"updateStatus"`
	UpdateTime       time.Time    `json:"updateTime,format:unixmilli"`
	FailureCode      int          `json:"failureCode"`
	FailureReason    string       `json:"failureReason"`
}

// GetStrategyOpenOrderService -- GET /fapi/v3/strategyOpenOrder (USER_DATA)
//
// Either strategyId or clientStrategyId must be set (mutually exclusive).
type GetStrategyOpenOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetStrategyOpenOrderService(strategyType StrategyType) *GetStrategyOpenOrderService {
	return &GetStrategyOpenOrderService{c: c, params: map[string]string{"strategyType": string(strategyType)}}
}

func (s *GetStrategyOpenOrderService) SetStrategyId(id int64) *GetStrategyOpenOrderService {
	s.params["strategyId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetStrategyOpenOrderService) SetClientStrategyId(id string) *GetStrategyOpenOrderService {
	s.params["clientStrategyId"] = id
	return s
}

func (s *GetStrategyOpenOrderService) Do(ctx context.Context) (*StrategyOrder, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/strategyOpenOrder", s.params).WithSignature()
	return request.Do[StrategyOrder](req)
}

// GetStrategyHistoryOrderService -- GET /fapi/v3/strategyHistoryOrder (USER_DATA)
//
// Either strategyId or clientStrategyId must be set (mutually exclusive).
// Maximum lookback window is 90 days.
type GetStrategyHistoryOrderService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetStrategyHistoryOrderService(strategyType StrategyType) *GetStrategyHistoryOrderService {
	return &GetStrategyHistoryOrderService{c: c, params: map[string]string{"strategyType": string(strategyType)}}
}

func (s *GetStrategyHistoryOrderService) SetStrategyId(id int64) *GetStrategyHistoryOrderService {
	s.params["strategyId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetStrategyHistoryOrderService) SetClientStrategyId(id string) *GetStrategyHistoryOrderService {
	s.params["clientStrategyId"] = id
	return s
}

func (s *GetStrategyHistoryOrderService) SetStartTime(t time.Time) *GetStrategyHistoryOrderService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetStrategyHistoryOrderService) SetEndTime(t time.Time) *GetStrategyHistoryOrderService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetStrategyHistoryOrderService) SetLimit(limit int) *GetStrategyHistoryOrderService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetStrategyHistoryOrderService) Do(ctx context.Context) (*StrategyOrder, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/strategyHistoryOrder", s.params).WithSignature()
	return request.Do[StrategyOrder](req)
}

// StrategyOrder is the shape returned by the strategy open/history queries.
type StrategyOrder struct {
	StrategyId       int64                 `json:"strategyId"`
	ClientStrategyId string                `json:"clientStrategyId"`
	StrategyType     StrategyType          `json:"strategyType"`
	StrategyStatus   string                `json:"strategyStatus"`
	BookTime         time.Time             `json:"bookTime,format:unixmilli"`
	UpdateTime       time.Time             `json:"updateTime,format:unixmilli"`
	SubOrders        []StrategyOrderDetail `json:"subOrders"`
}

type StrategyOrderDetail struct {
	StrategyId     int64           `json:"strategyId"`
	OrderId        int64           `json:"orderId"`
	ClientOrderId  string          `json:"clientOrderId"`
	Status         OrderStatus     `json:"status"`
	StrategySubId  int             `json:"strategySubId"`
	FirstDrivenId  int64           `json:"firstDrivenId"`
	FirstDrivenOn  string          `json:"firstDrivenOn"`
	FirstTrigger   string          `json:"firstTrigger"`
	SecondDrivenId int64           `json:"secondDrivenId"`
	SecondDrivenOn string          `json:"secondDrivenOn"`
	SecondTrigger  string          `json:"secondTrigger"`
	SecurityType   string          `json:"securityType"`
	Symbol         string          `json:"symbol"`
	Side           OrderSide       `json:"side"`
	PositionSide   PositionSide    `json:"positionSide"`
	Type           OrderType       `json:"type"`
	TimeInForce    TimeInForce     `json:"timeInForce"`
	Quantity       decimal.Decimal `json:"quantity"`
	ReduceOnly     bool            `json:"reduceOnly"`
	ClosePosition  bool            `json:"closePosition"`
	Price          decimal.Decimal `json:"price"`
	AvgPrice       decimal.Decimal `json:"avgPrice"`
	PriceProtect   bool            `json:"priceProtect"`
	StopPrice      decimal.Decimal `json:"stopPrice"`
	ActivatePrice  decimal.Decimal `json:"activatePrice"`
	CallbackRate   decimal.Decimal `json:"callbackRate"`
	WorkingType    WorkingType     `json:"workingType"`
	TriggerTime    time.Time       `json:"triggerTime,format:unixmilli"`
}
