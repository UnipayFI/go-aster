package chain

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// WithdrawSignatureType selects how userSignature was produced.
type WithdrawSignatureType string

const (
	WithdrawSignatureEOA        WithdrawSignatureType = "EOA"
	WithdrawSignatureSafeWallet WithdrawSignatureType = "SafeWallet"
)

// The four withdraw endpoints are authorized by userSignature -- a signature
// the *user's own* wallet key produces over the withdrawal parameters -- and
// not by the signer/agent scheme the rest of the SDK uses. Their documented
// parameter tables carry no user/nonce/signer/signature fields, so these
// services send exactly what the caller supplies and do NOT call WithSignature.
//
// The docs do not publish the EIP-712 typed-data schema behind userSignature,
// so the SDK cannot build it for you: produce it with your wallet tooling and
// pass the hex (or, for SafeWallet, the comma-separated set) in.

// PerpWithdrawService -- POST /aster-chain/v3/perp/user-withdraw (WITHDRAW)
//
// Withdraws from the perp account to an on-chain EVM address.
type PerpWithdrawService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewPerpWithdrawService(asset string, chainID int64, amount, fee decimal.Decimal, receiver, userNonce, userSignature string) *PerpWithdrawService {
	return &PerpWithdrawService{c: c, params: map[string]string{
		"asset":         asset,
		"chainId":       strconv.FormatInt(chainID, 10),
		"amount":        amount.String(),
		"fee":           fee.String(),
		"receiver":      receiver,
		"userNonce":     userNonce,
		"userSignature": userSignature,
	}}
}

// SetSignatureType marks the withdrawal as Safe-wallet signed. Defaults to EOA.
func (s *PerpWithdrawService) SetSignatureType(t WithdrawSignatureType) *PerpWithdrawService {
	s.params["signatureType"] = string(t)
	return s
}

func (s *PerpWithdrawService) Do(ctx context.Context) (*WithdrawResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/perp/user-withdraw", s.params)
	return request.Do[WithdrawResult](req)
}

// PerpSolanaWithdrawService -- POST /aster-chain/v3/perp/user-solana-withdraw (WITHDRAW)
//
// Withdraws from the perp account to a Solana address. Solana withdrawals have
// no signatureType: the signature is always the user's own Ed25519 signature.
type PerpSolanaWithdrawService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewPerpSolanaWithdrawService(asset string, chainID int64, amount, fee decimal.Decimal, receiver, userNonce, userSignature string) *PerpSolanaWithdrawService {
	return &PerpSolanaWithdrawService{c: c, params: map[string]string{
		"asset":         asset,
		"chainId":       strconv.FormatInt(chainID, 10),
		"amount":        amount.String(),
		"fee":           fee.String(),
		"receiver":      receiver,
		"userNonce":     userNonce,
		"userSignature": userSignature,
	}}
}

func (s *PerpSolanaWithdrawService) Do(ctx context.Context) (*WithdrawResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/perp/user-solana-withdraw", s.params)
	return request.Do[WithdrawResult](req)
}

// SpotWithdrawService -- POST /aster-chain/v3/spot/user-withdraw (WITHDRAW)
//
// Withdraws from the spot account to an on-chain EVM address. chainID is the
// chain the asset lives on; when the signature was produced against a
// different EIP-712 domain chainId, declare it with SetSignatureChainId.
type SpotWithdrawService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewSpotWithdrawService(asset string, chainID int64, amount, fee decimal.Decimal, receiver, userNonce, userSignature string) *SpotWithdrawService {
	return &SpotWithdrawService{c: c, params: map[string]string{
		"asset":         asset,
		"chainId":       strconv.FormatInt(chainID, 10),
		"amount":        amount.String(),
		"fee":           fee.String(),
		"receiver":      receiver,
		"userNonce":     userNonce,
		"userSignature": userSignature,
	}}
}

// SetSignatureChainId sets the chainId of the EIP-712 domain userSignature was
// produced under. Defaults to the asset's chainId.
func (s *SpotWithdrawService) SetSignatureChainId(chainID int64) *SpotWithdrawService {
	s.params["signatureChainId"] = strconv.FormatInt(chainID, 10)
	return s
}

// SetSignatureType marks the withdrawal as Safe-wallet signed. Defaults to EOA.
func (s *SpotWithdrawService) SetSignatureType(t WithdrawSignatureType) *SpotWithdrawService {
	s.params["signatureType"] = string(t)
	return s
}

func (s *SpotWithdrawService) Do(ctx context.Context) (*WithdrawResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/spot/user-withdraw", s.params)
	return request.Do[WithdrawResult](req)
}

// SpotSolanaWithdrawService -- POST /aster-chain/v3/spot/user-solana-withdraw (WITHDRAW)
//
// Withdraws from the spot account to a Solana address.
type SpotSolanaWithdrawService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewSpotSolanaWithdrawService(asset string, chainID int64, amount, fee decimal.Decimal, receiver, userNonce, userSignature string) *SpotSolanaWithdrawService {
	return &SpotSolanaWithdrawService{c: c, params: map[string]string{
		"asset":         asset,
		"chainId":       strconv.FormatInt(chainID, 10),
		"amount":        amount.String(),
		"fee":           fee.String(),
		"receiver":      receiver,
		"userNonce":     userNonce,
		"userSignature": userSignature,
	}}
}

func (s *SpotSolanaWithdrawService) Do(ctx context.Context) (*WithdrawResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/spot/user-solana-withdraw", s.params)
	return request.Do[WithdrawResult](req)
}

// WithdrawResult identifies an accepted withdrawal. hash is the on-chain
// transaction hash once the withdrawal has been broadcast.
type WithdrawResult struct {
	WithdrawID string `json:"withdrawId"`
	Hash       string `json:"hash"`
}

// GetPerpWithdrawInfoService -- GET /aster-chain/v3/perp/user-withdraw-info (USER_DATA)
//
// Returns the caller's withdrawal limits and withdrawable balances per asset
// and per chain. The doc lists no parameters, but the endpoint is USER_DATA and
// resolves "the current user", so the request is signed with the regular
// configured signer.
type GetPerpWithdrawInfoService struct {
	c *ChainClient
}

func (c *ChainClient) NewGetPerpWithdrawInfoService() *GetPerpWithdrawInfoService {
	return &GetPerpWithdrawInfoService{c: c}
}

func (s *GetPerpWithdrawInfoService) Do(ctx context.Context) (*WithdrawInfo, error) {
	req := request.Get(ctx, s.c, "/aster-chain/v3/perp/user-withdraw-info").WithSignature()
	return request.Do[WithdrawInfo](req)
}

// WithdrawInfo reports both the user's own daily allowance and the
// platform-wide one; a withdrawal is capped by whichever runs out first.
// Balances is keyed by asset name.
type WithdrawInfo struct {
	UserDailyLimit           decimal.Decimal                 `json:"userDailyLimit"`
	UserRemainingDailyLimit  decimal.Decimal                 `json:"userRemainingDailyLimit"`
	TotalDailyLimit          decimal.Decimal                 `json:"totalDailyLimit"`
	TotalRemainingDailyLimit decimal.Decimal                 `json:"totalRemainingDailyLimit"`
	Balances                 map[string]WithdrawAssetBalance `json:"balances"`
}

// WithdrawAssetBalance is one asset's withdrawal state. ChainBalances is keyed
// by chainId rendered as a decimal string.
type WithdrawAssetBalance struct {
	Currency                string                          `json:"currency"`
	SpotTotalWithdrawAmount decimal.Decimal                 `json:"spotTotalWithdrawAmount"`
	PerpTotalWithdrawAmount decimal.Decimal                 `json:"perpTotalWithdrawAmount"`
	DailyLimit              decimal.Decimal                 `json:"dailyLimit"`
	ChainBalances           map[string]WithdrawChainBalance `json:"chainBalances"`
}

// WithdrawChainBalance is one asset's withdrawal state on a single chain.
type WithdrawChainBalance struct {
	ChainID               int64           `json:"chainId"`
	SpotMaxWithdrawAmount decimal.Decimal `json:"spotMaxWithdrawAmount"`
	PerpMaxWithdrawAmount decimal.Decimal `json:"perpMaxWithdrawAmount"`
	ChainLimit            decimal.Decimal `json:"chainLimit"`
	WithdrawFee           decimal.Decimal `json:"withdrawFee"`
}

// GetDepositWithdrawHistoryService -- GET /aster-chain/v3/perp/deposit-withdraw-history (USER_DATA)
//
// Returns the caller's perp-account deposit and withdrawal records.
type GetDepositWithdrawHistoryService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewGetDepositWithdrawHistoryService() *GetDepositWithdrawHistoryService {
	return &GetDepositWithdrawHistoryService{c: c, params: map[string]string{}}
}

// SetChainId restricts the result to a single chain.
func (s *GetDepositWithdrawHistoryService) SetChainId(chainID int64) *GetDepositWithdrawHistoryService {
	s.params["chainId"] = strconv.FormatInt(chainID, 10)
	return s
}

func (s *GetDepositWithdrawHistoryService) Do(ctx context.Context) ([]DepositWithdrawRecord, error) {
	req := request.Get(ctx, s.c, "/aster-chain/v3/perp/deposit-withdraw-history", s.params).WithSignature()
	resp, err := request.Do[[]DepositWithdrawRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DepositWithdrawRecord is one deposit or withdrawal. Type is DEPOSIT or
// WITHDRAW; AccountType names the account side, e.g. perp. Time is unix
// milliseconds -- see TimeAt for a time.Time.
type DepositWithdrawRecord struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Asset       string          `json:"asset"`
	Amount      decimal.Decimal `json:"amount"`
	State       string          `json:"state"`
	TxHash      string          `json:"txHash"`
	Time        int64           `json:"time"`
	ChainID     int64           `json:"chainId"`
	AccountType string          `json:"accountType"`
}

// TimeAt returns Time as a time.Time.
func (r DepositWithdrawRecord) TimeAt() time.Time {
	return time.UnixMilli(r.Time)
}

// EstimateWithdrawFeeService -- GET /aster-chain/v3/withdraw/estimateFee (NONE)
//
// Estimates the gas fee for withdrawing an asset on a chain. Public: no
// signature required.
type EstimateWithdrawFeeService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewEstimateWithdrawFeeService(chainID int64, asset string) *EstimateWithdrawFeeService {
	return &EstimateWithdrawFeeService{c: c, params: map[string]string{
		"chainId": strconv.FormatInt(chainID, 10),
		"asset":   asset,
	}}
}

func (s *EstimateWithdrawFeeService) Do(ctx context.Context) (*WithdrawFeeEstimate, error) {
	req := request.Get(ctx, s.c, "/aster-chain/v3/withdraw/estimateFee", s.params)
	return request.Do[WithdrawFeeEstimate](req)
}

// WithdrawFeeEstimate breaks the fee down into its gas inputs and their fiat
// value. GasPrice and GasLimit arrive as JSON numbers; the prices and costs
// arrive as strings.
type WithdrawFeeEstimate struct {
	GasPrice    int64           `json:"gasPrice"`
	GasLimit    int64           `json:"gasLimit"`
	NativePrice decimal.Decimal `json:"nativePrice"`
	TokenPrice  decimal.Decimal `json:"tokenPrice"`
	GasCost     decimal.Decimal `json:"gasCost"`
	GasUsdValue decimal.Decimal `json:"gasUsdValue"`
}
