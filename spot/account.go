package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// GetAccountService -- GET /api/v3/account (USER_DATA)
type GetAccountService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

func (s *GetAccountService) Do(ctx context.Context) (*AccountInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/account").WithSignature()
	return request.Do[AccountInfo](req)
}

type AccountInfo struct {
	FeeTier      int64     `json:"feeTier"`
	CanTrade     bool      `json:"canTrade"`
	CanDeposit   bool      `json:"canDeposit"`
	CanWithdraw  bool      `json:"canWithdraw"`
	CanBurnAsset bool      `json:"canBurnAsset"`
	UpdateTime   int64     `json:"updateTime"`
	Balances     []Balance `json:"balances"`
}

type Balance struct {
	Asset  string          `json:"asset"`
	Free   decimal.Decimal `json:"free"`
	Locked decimal.Decimal `json:"locked"`
}

// GetUserTradesService -- GET /api/v3/userTrades (USER_DATA)
type GetUserTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetUserTradesService() *GetUserTradesService {
	return &GetUserTradesService{c: c, params: map[string]string{}}
}

func (s *GetUserTradesService) SetSymbol(symbol string) *GetUserTradesService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetUserTradesService) SetOrderId(id int64) *GetUserTradesService {
	s.params["orderId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetUserTradesService) SetStartTime(t time.Time) *GetUserTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetEndTime(t time.Time) *GetUserTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetFromId(id int64) *GetUserTradesService {
	s.params["fromId"] = strconv.FormatInt(id, 10)
	return s
}

func (s *GetUserTradesService) SetLimit(limit int) *GetUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetUserTradesService) Do(ctx context.Context) ([]UserTrade, error) {
	req := request.Get(ctx, s.c, "/api/v3/userTrades", s.params).WithSignature()
	resp, err := request.Do[[]UserTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type UserTrade struct {
	Symbol          string          `json:"symbol"`
	ID              int64           `json:"id"`
	OrderID         int64           `json:"orderId"`
	Side            OrderSide       `json:"side"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	QuoteQty        decimal.Decimal `json:"quoteQty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	Time            time.Time       `json:"time,format:unixmilli"`
	CounterpartyID  int64           `json:"counterpartyId"`
	CreateUpdateID  *int64          `json:"createUpdateId"`
	Maker           bool            `json:"maker"`
	Buyer           bool            `json:"buyer"`
}

// PerpSpotTransferService -- POST /api/v3/asset/wallet/transfer (TRADE)
//
// Move funds between the spot wallet and the perpetual (futures) wallet.
type PerpSpotTransferService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewPerpSpotTransferService(asset string, amount decimal.Decimal, kind TransferKindType, clientTranID string) *PerpSpotTransferService {
	return &PerpSpotTransferService{c: c, params: map[string]string{
		"amount":       amount.String(),
		"asset":        asset,
		"clientTranId": clientTranID,
		"kindType":     string(kind),
	}}
}

func (s *PerpSpotTransferService) Do(ctx context.Context) (*TransferResponse, error) {
	req := request.Post(s.c, ctx, "/api/v3/asset/wallet/transfer", s.params).WithSignature()
	return request.Do[TransferResponse](req)
}

type TransferResponse struct {
	TranID int64  `json:"tranId"`
	Status string `json:"status"`
}

// GetWithdrawFeeService -- GET /api/v3/aster/withdraw/estimateFee (NONE)
type GetWithdrawFeeService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetWithdrawFeeService(chainID string, asset string) *GetWithdrawFeeService {
	return &GetWithdrawFeeService{c: c, params: map[string]string{
		"chainId": chainID,
		"asset":   asset,
	}}
}

func (s *GetWithdrawFeeService) Do(ctx context.Context) (*WithdrawFeeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/aster/withdraw/estimateFee", s.params)
	return request.Do[WithdrawFeeResponse](req)
}

type WithdrawFeeResponse struct {
	TokenPrice  decimal.Decimal `json:"tokenPrice"`
	GasCost     decimal.Decimal `json:"gasCost"`
	GasUsdValue decimal.Decimal `json:"gasUsdValue"`
}

// WithdrawService -- POST /api/v3/aster/user-withdraw (USER_DATA)
//
// Withdraws require a separately-signed `userSignature` parameter generated
// from a different EIP-712 typed-data structure (the "Action/Withdraw" type
// with verifyingContract=ZeroAddress, name="Aster"). Because the typed-data
// schema is unrelated to the request-signing flow, the caller must compute
// userSignature themselves -- see the doc snippet at
// /tmp/api-docs/V3(Recommended)/EN/aster-finance-spot-api-v3.md lines
// 1430-1465 for the exact field set. The SDK only attaches the standard V3
// request signature.
type WithdrawService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewWithdrawService(chainID, asset, amount, fee, receiver, nonce, userSignature string) *WithdrawService {
	return &WithdrawService{c: c, params: map[string]string{
		"chainId":       chainID,
		"asset":         asset,
		"amount":        amount,
		"fee":           fee,
		"receiver":      receiver,
		"nonce":         nonce,
		"userSignature": userSignature,
	}}
}

func (s *WithdrawService) Do(ctx context.Context) (*WithdrawResponse, error) {
	req := request.Post(s.c, ctx, "/api/v3/aster/user-withdraw", s.params).WithSignature()
	return request.Do[WithdrawResponse](req)
}

type WithdrawResponse struct {
	WithdrawID string `json:"withdrawId"`
	Hash       string `json:"hash"`
}
