package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/internal/request"
	"github.com/shopspring/decimal"
)

type GetAccountService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c, params: map[string]string{}}
}

func (s *GetAccountService) Do(ctx context.Context) (*AccountResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/account", s.params).Sign()
	return request.Do[AccountResponse](req)
}

type AccountResponse struct {
	FeeTier      int              `json:"feeTier"`
	CanTrade     bool             `json:"canTrade"`
	CanDeposit   bool             `json:"canDeposit"`
	CanWithdraw  bool             `json:"canWithdraw"`
	CanBurnAsset bool             `json:"canBurnAsset"`
	UpdateTime   time.Time        `json:"updateTime,format:unixmilli"`
	Balances     []AccountBalance `json:"balances"`
}

type AccountBalance struct {
	Asset  string          `json:"asset"`
	Free   decimal.Decimal `json:"free"`
	Locked decimal.Decimal `json:"locked"`
}

type GetUserTradesService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewGetUserTradesService() *GetUserTradesService {
	return &GetUserTradesService{c: c, params: map[string]string{}}
}

func (s *GetUserTradesService) SetSymbol(symbol string) *GetUserTradesService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetUserTradesService) SetOrderId(orderId int64) *GetUserTradesService {
	s.params["orderId"] = strconv.FormatInt(orderId, 10)
	return s
}

func (s *GetUserTradesService) SetStartTime(startTime time.Time) *GetUserTradesService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetEndTime(endTime time.Time) *GetUserTradesService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetUserTradesService) SetFromId(fromId int64) *GetUserTradesService {
	s.params["fromId"] = strconv.FormatInt(fromId, 10)
	return s
}

func (s *GetUserTradesService) SetLimit(limit int) *GetUserTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetUserTradesService) Do(ctx context.Context) ([]UserTradeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/userTrades", s.params).Sign()
	trades, err := request.Do[[]UserTradeResponse](req)
	if err != nil {
		return nil, err
	}
	return *trades, nil
}

type UserTradeResponse struct {
	Symbol          string          `json:"symbol"`
	Id              int64           `json:"id"`
	OrderId         int64           `json:"orderId"`
	Side            OrderSide       `json:"side"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	QuoteQty        decimal.Decimal `json:"quoteQty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	Time            time.Time       `json:"time,format:unixmilli"`
	CounterpartyId  int64           `json:"counterpartyId"`
	CreateUpdateId  *int64          `json:"createUpdateId"`
	Maker           bool            `json:"maker"`
	Buyer           bool            `json:"buyer"`
}

type WalletTransferService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewWalletTransferService(amount float64, asset string, clientTranId string, kindType TransferType) *WalletTransferService {
	return &WalletTransferService{
		c: c,
		params: map[string]string{
			"amount":       strconv.FormatFloat(amount, 'f', -1, 64),
			"asset":        asset,
			"clientTranId": clientTranId,
			"kindType":     string(kindType),
		},
	}
}

func (s *WalletTransferService) Do(ctx context.Context) (*TransferResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/asset/wallet/transfer", s.params).Sign()
	return request.Do[TransferResponse](req)
}

type TransferResponse struct {
	TranId int64  `json:"tranId"`
	Status string `json:"status"`
}

type SendToAddressService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewSendToAddressService(amount float64, asset string, toAddress string) *SendToAddressService {
	return &SendToAddressService{
		c: c,
		params: map[string]string{
			"amount":    strconv.FormatFloat(amount, 'f', -1, 64),
			"asset":     asset,
			"toAddress": toAddress,
		},
	}
}

func (s *SendToAddressService) SetClientTranId(clientTranId string) *SendToAddressService {
	s.params["clientTranId"] = clientTranId
	return s
}

func (s *SendToAddressService) Do(ctx context.Context) (*TransferResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/asset/sendToAddress", s.params).Sign()
	return request.Do[TransferResponse](req)
}

type GetWithdrawFeeService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewGetWithdrawFeeService(chainId string, asset string) *GetWithdrawFeeService {
	return &GetWithdrawFeeService{
		c: c,
		params: map[string]string{
			"chainId": chainId,
			"asset":   asset,
		},
	}
}

func (s *GetWithdrawFeeService) Do(ctx context.Context) (*WithdrawFeeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v1/aster/withdraw/estimateFee", s.params)
	return request.Do[WithdrawFeeResponse](req)
}

type WithdrawFeeResponse struct {
	TokenPrice  decimal.Decimal `json:"tokenPrice"`
	GasCost     decimal.Decimal `json:"gasCost"`
	GasUsdValue decimal.Decimal `json:"gasUsdValue"`
}

type WithdrawService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewWithdrawService(chainId string, asset string, amount string, fee string, receiver string, nonce string, userSignature string) *WithdrawService {
	return &WithdrawService{
		c: c,
		params: map[string]string{
			"chainId":       chainId,
			"asset":         asset,
			"amount":        amount,
			"fee":           fee,
			"receiver":      receiver,
			"nonce":         nonce,
			"userSignature": userSignature,
		},
	}
}

func (s *WithdrawService) Do(ctx context.Context) (*WithdrawResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/aster/user-withdraw", s.params).Sign()
	return request.Do[WithdrawResponse](req)
}

type WithdrawResponse struct {
	WithdrawId string `json:"withdrawId"`
	Hash       string `json:"hash"`
}

type GetNonceService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewGetNonceService(address string, userOperationType string) *GetNonceService {
	return &GetNonceService{
		c: c,
		params: map[string]string{
			"address":           address,
			"userOperationType": userOperationType,
		},
	}
}

func (s *GetNonceService) SetNetwork(network string) *GetNonceService {
	s.params["network"] = network
	return s
}

func (s *GetNonceService) Do(ctx context.Context) (*int64, error) {
	req := request.Post(s.c, ctx, "/api/v1/getNonce", s.params)
	return request.Do[int64](req)
}

type CreateApiKeyService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewCreateApiKeyService(address string, userOperationType string, userSignature string, desc string) *CreateApiKeyService {
	return &CreateApiKeyService{
		c: c,
		params: map[string]string{
			"address":           address,
			"userOperationType": userOperationType,
			"userSignature":     userSignature,
			"desc":              desc,
		},
	}
}

func (s *CreateApiKeyService) SetNetwork(network string) *CreateApiKeyService {
	s.params["network"] = network
	return s
}

func (s *CreateApiKeyService) SetApikeyIP(apikeyIP string) *CreateApiKeyService {
	s.params["apikeyIP"] = apikeyIP
	return s
}

func (s *CreateApiKeyService) Do(ctx context.Context) (*CreateApiKeyResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/createApiKey", s.params)
	return request.Do[CreateApiKeyResponse](req)
}

type CreateApiKeyResponse struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type CreateListenKeyService struct {
	c *SpotClient
}

func (c *SpotClient) NewCreateListenKeyService() *CreateListenKeyService {
	return &CreateListenKeyService{c: c}
}

func (s *CreateListenKeyService) Do(ctx context.Context) (*ListenKeyResponse, error) {
	req := request.Post(s.c, ctx, "/api/v1/listenKey").SetApiKeyHeader()
	return request.Do[ListenKeyResponse](req)
}

type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

type ExtendListenKeyService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewExtendListenKeyService(listenKey string) *ExtendListenKeyService {
	return &ExtendListenKeyService{
		c:      c,
		params: map[string]string{"listenKey": listenKey},
	}
}

func (s *ExtendListenKeyService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/api/v1/listenKey", s.params).SetApiKeyHeader()
	_, err := request.Do[struct{}](req)
	return err
}

type CloseListenKeyService struct {
	c *SpotClient

	params map[string]string
}

func (c *SpotClient) NewCloseListenKeyService(listenKey string) *CloseListenKeyService {
	return &CloseListenKeyService{
		c:      c,
		params: map[string]string{"listenKey": listenKey},
	}
}

func (s *CloseListenKeyService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v1/listenKey", s.params).SetApiKeyHeader()
	_, err := request.Do[struct{}](req)
	return err
}
