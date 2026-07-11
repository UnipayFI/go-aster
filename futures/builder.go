package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// GetBuilderUserTradesService -- GET /fapi/v3/builder/userTrades (USER_DATA)
//
// Queries the paginated trade history of users trading under the caller's
// builder code. The authenticated account is the builder identity, so there is
// no separate builder address parameter. If neither startTime nor endTime is
// set, the recent 7 days are returned; the resolved startTime may not be more
// than 30 days before now, and endTime may not be more than 1 day ahead of
// server time. Results are ordered by trade time descending.
type GetBuilderUserTradesService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderUserTradesService() *GetBuilderUserTradesService {
	return &GetBuilderUserTradesService{c: c, params: map[string]string{}}
}

func (s *GetBuilderUserTradesService) SetStartTime(t time.Time) *GetBuilderUserTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderUserTradesService) SetEndTime(t time.Time) *GetBuilderUserTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderUserTradesService) SetPage(page int) *GetBuilderUserTradesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *GetBuilderUserTradesService) SetLimit(limit int) *GetBuilderUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBuilderUserTradesService) Do(ctx context.Context) (*BuilderUserTrades, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/userTrades", s.params).WithSignature()
	return request.Do[BuilderUserTrades](req)
}

// BuilderUserTrades is the paginated response of GetBuilderUserTradesService.
type BuilderUserTrades struct {
	Total       int64              `json:"total"`
	CurrentPage int                `json:"currentPage"`
	TotalPages  int                `json:"totalPages"`
	PageSize    int                `json:"pageSize"`
	HasMore     bool               `json:"hasMore"`
	Rows        []BuilderUserTrade `json:"rows"`
}

// BuilderUserTrade is one trade executed by a user under the caller's builder
// code. baseType is CRYPTO or STOCK; totalQuota is the notional (price*qty).
type BuilderUserTrade struct {
	TradeID        int64           `json:"tradeId"`
	InsertTime     time.Time       `json:"insertTime,format:unixmilli"`
	Symbol         string          `json:"symbol"`
	PositionSide   PositionSide    `json:"positionSide"`
	Price          decimal.Decimal `json:"price"`
	Qty            decimal.Decimal `json:"qty"`
	BaseAsset      string          `json:"baseAsset"`
	BaseType       string          `json:"baseType"`
	QuoteAsset     string          `json:"quoteAsset"`
	Side           OrderSide       `json:"side"`
	ActiveBuy      bool            `json:"activeBuy"`
	FeeAsset       string          `json:"feeAsset"`
	TotalQuota     decimal.Decimal `json:"totalQuota"`
	Fee            decimal.Decimal `json:"fee"`
	OrderID        int64           `json:"orderId"`
	RealizedProfit decimal.Decimal `json:"realizedProfit"`
	MarginAsset    string          `json:"marginAsset"`
	UserAddress    string          `json:"userAddress"`
	BuilderFee     decimal.Decimal `json:"builderFee"`
}

// GetBuilderApprovedUserListService -- GET /fapi/v3/builder/approvedUserList (USER_DATA)
//
// Queries the users who have approved the caller's address as their builder.
// The authenticated account is the builder identity; there is no separate
// builder address parameter. Returns at most 1000 records.
type GetBuilderApprovedUserListService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetBuilderApprovedUserListService() *GetBuilderApprovedUserListService {
	return &GetBuilderApprovedUserListService{c: c, params: map[string]string{}}
}

// SetStartTime only returns users approved after this timestamp. When unset the
// server applies no time filter and returns all approved users.
func (s *GetBuilderApprovedUserListService) SetStartTime(t time.Time) *GetBuilderApprovedUserListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBuilderApprovedUserListService) Do(ctx context.Context) ([]BuilderApprovedUser, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/builder/approvedUserList", s.params).WithSignature()
	resp, err := request.Do[[]BuilderApprovedUser](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BuilderApprovedUser is one user who approved the caller as their builder.
type BuilderApprovedUser struct {
	UserAddress    string          `json:"userAddress"`
	BuilderAddress string          `json:"builderAddress"`
	MaxFeeRate     decimal.Decimal `json:"maxFeeRate"`
	BuilderName    string          `json:"builderName"`
	ApproveTime    time.Time       `json:"approveTime,format:unixmilli"`
}
