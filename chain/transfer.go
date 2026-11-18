package chain

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// TransferKindType is the direction of a wallet transfer between the spot
// wallet and the perp account. Both wallet/transfer endpoints move funds
// between the same two accounts, so both directions are valid on either one.
type TransferKindType string

const (
	TransferSpotToPerp TransferKindType = "SPOT_FUTURE"
	TransferPerpToSpot TransferKindType = "FUTURE_SPOT"
)

// Like the futures sub-account flows, the wallet/transfer endpoints are signed
// with the *user's own* wallet private key rather than the configured
// signer/agent key. Master and source-account keys usually live in cold
// storage, so these services take an already-computed signature and send the
// payload as-is -- they do NOT call WithSignature. Build the message body from
// the parameters below and sign it with request.SignEIP712V3.

// PerpWalletTransferService -- POST /aster-chain/v3/perp/wallet/transfer (TRADE)
//
// Moves an asset between the spot wallet and the perp account, initiated from
// the perp side. `user` is the source account's wallet address and must be the
// address whose key produced the signature.
type PerpWalletTransferService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewPerpWalletTransferService(asset string, amount decimal.Decimal, clientTranID string, kind TransferKindType, user string, nonce int64, signature string) *PerpWalletTransferService {
	return &PerpWalletTransferService{c: c, params: walletTransferParams(asset, amount, clientTranID, kind, user, nonce, signature)}
}

func (s *PerpWalletTransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/perp/wallet/transfer", s.params)
	return request.Do[TransferResult](req)
}

// SpotWalletTransferService -- POST /aster-chain/v3/spot/wallet/transfer (TRADE)
//
// The spot-side counterpart of PerpWalletTransferService: same parameters,
// same response.
type SpotWalletTransferService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewSpotWalletTransferService(asset string, amount decimal.Decimal, clientTranID string, kind TransferKindType, user string, nonce int64, signature string) *SpotWalletTransferService {
	return &SpotWalletTransferService{c: c, params: walletTransferParams(asset, amount, clientTranID, kind, user, nonce, signature)}
}

func (s *SpotWalletTransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(s.c, ctx, "/aster-chain/v3/spot/wallet/transfer", s.params)
	return request.Do[TransferResult](req)
}

func walletTransferParams(asset string, amount decimal.Decimal, clientTranID string, kind TransferKindType, user string, nonce int64, signature string) map[string]string {
	return map[string]string{
		"asset":        asset,
		"amount":       amount.String(),
		"clientTranId": clientTranID,
		"kindType":     string(kind),
		"nonce":        strconv.FormatInt(nonce, 10),
		"user":         user,
		"signature":    signature,
	}
}

type TransferResult struct {
	TranID int64  `json:"tranId"`
	Status string `json:"status"`
}

// BatchTransferStatus reports the processing state of a batch transfer or one
// of its items.
type BatchTransferStatus string

const (
	BatchTransferPending    BatchTransferStatus = "I"
	BatchTransferProcessing BatchTransferStatus = "P"
	BatchTransferSuccess    BatchTransferStatus = "S"
	BatchTransferFailed     BatchTransferStatus = "F"
	BatchTransferPartial    BatchTransferStatus = "M"
)

// BatchTransferItem is one recipient entry in BatchTransferService. Amount is
// kept as a string (greater than 0, at most 8 decimal places) so the marshalled
// JSON bytes match what the caller signed.
type BatchTransferItem struct {
	ToAddress       string `json:"toAddress"`
	Asset           string `json:"asset"`
	Amount          string `json:"amount"`
	TransferChainID *int64 `json:"transferChainId,omitempty"`
}

// BatchTransferService -- POST /aster-chain/v3/batch-transfer (WITHDRAW)
//
// Transfers assets from the current account to up to 10 registered Aster-Chain
// addresses in one asynchronous batch. Not available while the account is in
// PRIVATE mode. Any invalid item (bad, own or unregistered address, duplicate
// toAddress+asset pair, or amount <= 0 / > 8 decimals) rejects the whole batch.
// clientTranID (20–64 characters) makes the call idempotent.
//
// Signed with the `user` account's own wallet private key rather than the
// configured signer key, so the caller supplies the signature and this service
// does NOT call WithSignature. Sign the URL-encoded message body via
// request.SignEIP712V3 (where itemsJSON is json.Marshal(items)):
//
//	clientTranId={clientTranId}&items={url.QueryEscape(itemsJSON)}&nonce={nonce}&user={user}
type BatchTransferService struct {
	c            *ChainClient
	clientTranID string
	items        []BatchTransferItem
	user         string
	nonce        int64
	signature    string
}

func (c *ChainClient) NewBatchTransferService(clientTranID string, items []BatchTransferItem, user string, nonce int64, signature string) *BatchTransferService {
	return &BatchTransferService{
		c:            c,
		clientTranID: clientTranID,
		items:        items,
		user:         user,
		nonce:        nonce,
		signature:    signature,
	}
}

func (s *BatchTransferService) Do(ctx context.Context) (*BatchTransferResult, error) {
	itemsJSON, err := json.Marshal(s.items)
	if err != nil {
		return nil, err
	}
	req := request.Post(s.c, ctx, "/aster-chain/v3/batch-transfer", map[string]string{
		"clientTranId": s.clientTranID,
		"items":        string(itemsJSON),
		"nonce":        strconv.FormatInt(s.nonce, 10),
		"user":         s.user,
		"signature":    s.signature,
	})
	return request.Do[BatchTransferResult](req)
}

// BatchTransferResult confirms that a batch transfer has been accepted; query
// GetBatchTransferHistoryService for final per-item results.
type BatchTransferResult struct {
	ClientTranID string                      `json:"clientTranId"`
	Status       BatchTransferStatus         `json:"status"`
	TotalCount   int                         `json:"totalCount"`
	Items        []BatchTransferAcceptedItem `json:"items"`
}

type BatchTransferAcceptedItem struct {
	ToAddress  string `json:"toAddress"`
	Asset      string `json:"asset"`
	TransferID string `json:"transferId"`
}

// GetBatchTransferHistoryService -- GET /aster-chain/v3/batch-transfer/history (USER_DATA)
//
// Queries batch transfers submitted by the authenticated account along with
// per-item results. Passing only SetClientTranId (without startTime or endTime)
// looks up that batch regardless of creation time; otherwise batches are
// filtered by creation time (startTime inclusive, endTime exclusive; default
// last 7 days, maximum span 31 days).
type GetBatchTransferHistoryService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewGetBatchTransferHistoryService() *GetBatchTransferHistoryService {
	return &GetBatchTransferHistoryService{c: c, params: map[string]string{}}
}

func (s *GetBatchTransferHistoryService) SetClientTranId(id string) *GetBatchTransferHistoryService {
	s.params["clientTranId"] = id
	return s
}

func (s *GetBatchTransferHistoryService) SetStatus(status BatchTransferStatus) *GetBatchTransferHistoryService {
	s.params["status"] = string(status)
	return s
}

func (s *GetBatchTransferHistoryService) SetStartTime(t time.Time) *GetBatchTransferHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBatchTransferHistoryService) SetEndTime(t time.Time) *GetBatchTransferHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPage sets the 1-based page number (default 1).
func (s *GetBatchTransferHistoryService) SetPage(page int) *GetBatchTransferHistoryService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the number of batches per page (default 20, maximum 100).
func (s *GetBatchTransferHistoryService) SetLimit(limit int) *GetBatchTransferHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBatchTransferHistoryService) Do(ctx context.Context) (*BatchTransferHistory, error) {
	req := request.Get(ctx, s.c, "/aster-chain/v3/batch-transfer/history", s.params).WithSignature()
	return request.Do[BatchTransferHistory](req)
}

type BatchTransferHistory struct {
	Total int64                 `json:"total"`
	Rows  []BatchTransferRecord `json:"rows"`
}

type BatchTransferRecord struct {
	ClientTranID    string                `json:"clientTranId"`
	Status          BatchTransferStatus   `json:"status"`
	CreateTime      time.Time             `json:"createTime,format:unixmilli"`
	UpdateTime      time.Time             `json:"updateTime,format:unixmilli"`
	TotalCount      int                   `json:"totalCount"`
	SuccessCount    int                   `json:"successCount"`
	FailCount       int                   `json:"failCount"`
	PartialCount    int                   `json:"partialCount"`
	ProcessingCount int                   `json:"processingCount"`
	FailedItems     []BatchTransferDetail `json:"failedItems"`
	Details         []BatchTransferDetail `json:"details"`
}

type BatchTransferDetail struct {
	TransferID      string              `json:"transferId"`
	ToAddress       string              `json:"toAddress"`
	Asset           string              `json:"asset"`
	Amount          decimal.Decimal     `json:"amount"`
	TransferChainID *int64              `json:"transferChainId"`
	SuccessAmount   *decimal.Decimal    `json:"successAmount"`
	Status          BatchTransferStatus `json:"status"`
	ErrorCode       *string             `json:"errorCode"`
	ErrorMsg        *string             `json:"errorMsg"`
	CreateTime      time.Time           `json:"createTime,format:unixmilli"`
	UpdateTime      time.Time           `json:"updateTime,format:unixmilli"`
}
