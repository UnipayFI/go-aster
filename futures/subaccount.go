package futures

import (
	"context"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// Sub-account endpoints are unusual: they require signatures from the *master*
// (and sometimes the *child*) wallet's own private key, NOT the signer/agent
// private key configured on the client. Because master keys are usually held
// in cold storage and should never live inside an SDK process, these services
// accept already-computed signatures as inputs. The caller is expected to:
//
//  1. Build the exact message-body string per the V3 docs (see comments on
//     each service for the format).
//  2. Sign that string via request.SignEIP712V3 (chainId=1666 by default for
//     mainnet) using the appropriate wallet private key.
//  3. Pass the resulting hex signature(s) to the service.
//
// Because every parameter is sent as-is (no automatic nonce/signer injection),
// these services do NOT call WithSignature -- the request payload is already
// fully signed.

// BindSubAccountService -- POST /fapi/v3/sub-accounts/bind (USER_DATA)
//
// Master + child two-signature flow.
//
// Child message body: childAddress={childAddress}&name={name}&nonce={nonce}&user={user}
// Master message body (with childSignature appended):
//
//	childAddress={childAddress}&name={name}&nonce={nonce}&user={user}&childSignature={childSignature}
type BindSubAccountService struct {
	c              *FuturesClient
	childAddress   string
	name           string
	nonce          int64
	user           string
	childSignature string
	signature      string
}

func (c *FuturesClient) NewBindSubAccountService(childAddress, name, user string, nonce int64, childSignature, masterSignature string) *BindSubAccountService {
	return &BindSubAccountService{
		c:              c,
		childAddress:   childAddress,
		name:           name,
		nonce:          nonce,
		user:           user,
		childSignature: childSignature,
		signature:      masterSignature,
	}
}

func (s *BindSubAccountService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/sub-accounts/bind", map[string]string{
		"childAddress":   s.childAddress,
		"name":           s.name,
		"nonce":          formatInt64(s.nonce),
		"user":           s.user,
		"childSignature": s.childSignature,
		"signature":      s.signature,
	})
	return request.Do[GenericCodeMsg](req)
}

// CreateSubAccountService -- POST /fapi/v3/createSubAccount (TRADE)
//
// Child message body:
//
//	subAccountName={subAccountName}&subSourceAddr={subSourceAddr}&nonce={nonce}&user={user}&signer={signer}
//
// Master message body (with childSignature appended):
//
//	subAccountName={subAccountName}&subSourceAddr={subSourceAddr}&nonce={nonce}&user={user}&signer={signer}&childSignature={childSignature}
type CreateSubAccountService struct {
	c              *FuturesClient
	subSourceAddr  string
	subAccountName string
	nonce          int64
	user           string
	signer         string
	childSignature string
	signature      string
}

func (c *FuturesClient) NewCreateSubAccountService(subSourceAddr, subAccountName, user, signer string, nonce int64, childSignature, masterSignature string) *CreateSubAccountService {
	return &CreateSubAccountService{
		c:              c,
		subSourceAddr:  subSourceAddr,
		subAccountName: subAccountName,
		nonce:          nonce,
		user:           user,
		signer:         signer,
		childSignature: childSignature,
		signature:      masterSignature,
	}
}

func (s *CreateSubAccountService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/createSubAccount", map[string]string{
		"subSourceAddr":  s.subSourceAddr,
		"subAccountName": s.subAccountName,
		"nonce":          formatInt64(s.nonce),
		"user":           s.user,
		"signer":         s.signer,
		"childSignature": s.childSignature,
		"signature":      s.signature,
	})
	return request.Do[GenericCodeMsg](req)
}

// GetSubAccountListService -- GET /fapi/v3/getSubAccountList (USER_DATA)
//
// Signs with the *signer* private key (EIP-712), unique among sub-account
// endpoints. The SDK can build this signature for you because it uses the
// regular configured signer, so this service does call WithSignature.
//
// Note: the doc explicitly states this endpoint always uses the EIP-712
// signer scheme even when the master account is Solana. WithSignature here
// produces exactly the right payload (nonce/signer/signature appended).
type GetSubAccountListService struct {
	c *FuturesClient
}

func (c *FuturesClient) NewGetSubAccountListService() *GetSubAccountListService {
	return &GetSubAccountListService{c: c}
}

func (s *GetSubAccountListService) Do(ctx context.Context) ([]SubAccountListEntry, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/getSubAccountList").WithSignature()
	resp, err := request.Do[[]SubAccountListEntry](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SubAccountListEntry struct {
	AccountID      int64  `json:"accountId"`
	SubAccountName string `json:"subAccountName"`
	ParentAccount  bool   `json:"parentAccount"`
}

// UpdateSubAccountService -- POST /fapi/v3/updateSubAccount (TRADE)
//
// Master signature only. Message body (omitting nil optional params):
//
//	subSourceAddr={subSourceAddr}&nonce={nonce}&user={user}&signer={signer}[&subAccountName={...}][&status={...}]
//
// At least one of subAccountName / status must be provided.
type UpdateSubAccountService struct {
	c              *FuturesClient
	subSourceAddr  string
	nonce          int64
	user           string
	signer         string
	subAccountName string
	status         string // NORMAL | FROZEN
	signature      string
}

func (c *FuturesClient) NewUpdateSubAccountService(subSourceAddr, user, signer string, nonce int64, masterSignature string) *UpdateSubAccountService {
	return &UpdateSubAccountService{
		c:             c,
		subSourceAddr: subSourceAddr,
		nonce:         nonce,
		user:          user,
		signer:        signer,
		signature:     masterSignature,
	}
}

func (s *UpdateSubAccountService) SetSubAccountName(name string) *UpdateSubAccountService {
	s.subAccountName = name
	return s
}

func (s *UpdateSubAccountService) SetStatus(status string) *UpdateSubAccountService {
	s.status = status
	return s
}

func (s *UpdateSubAccountService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	params := map[string]string{
		"subSourceAddr": s.subSourceAddr,
		"nonce":         formatInt64(s.nonce),
		"user":          s.user,
		"signer":        s.signer,
		"signature":     s.signature,
	}
	if s.subAccountName != "" {
		params["subAccountName"] = s.subAccountName
	}
	if s.status != "" {
		params["status"] = s.status
	}
	req := request.Post(s.c, ctx, "/fapi/v3/updateSubAccount", params)
	return request.Do[GenericCodeMsg](req)
}

// SubAccountTransferKindType discriminates the four transfer flows.
type SubAccountTransferKindType string

const (
	SubAccountTransferFutureToFuture SubAccountTransferKindType = "FUTURE_FUTURE"
	SubAccountTransferFutureToSpot   SubAccountTransferKindType = "FUTURE_SPOT"
	SubAccountTransferSpotToFuture   SubAccountTransferKindType = "SPOT_FUTURE"
	SubAccountTransferSpotToSpot     SubAccountTransferKindType = "SPOT_SPOT"
)

// SubAccountTransferService -- POST /fapi/v3/subAccountTransfer (TRADE)
//
// Identify the signing account with either SetUser or SetSigner -- only one is
// needed, and the signature must be produced with the private key of whichever
// one is passed. `user` is the signing account address: the master account in
// most scenarios, or the sub-account address when a sub-account initiates the
// transfer itself. `signer` is an agent wallet associated with that user; it
// must already be registered and approved (see RegisterAndApproveAgentService)
// before it can authorize transfers on the user's behalf.
//
// Message body (a field absent from the request is also absent here):
//
// Without fromAccountAddress:
//
//	toAccountAddress={...}&asset={...}&amount={...}&kindType={...}&nonce={...}&user={...}&signer={...}
//
// With fromAccountAddress:
//
//	toAccountAddress={...}&asset={...}&amount={...}&kindType={...}&nonce={...}&user={...}&signer={...}&fromAccountAddress={...}
type SubAccountTransferService struct {
	c                  *FuturesClient
	toAccountAddress   string
	asset              string
	amount             string
	kindType           SubAccountTransferKindType
	nonce              int64
	user               string
	signer             string
	fromAccountAddress string
	signature          string
}

func (c *FuturesClient) NewSubAccountTransferService(toAddr, asset, amount string, kind SubAccountTransferKindType, nonce int64, signature string) *SubAccountTransferService {
	return &SubAccountTransferService{
		c:                c,
		toAccountAddress: toAddr,
		asset:            asset,
		amount:           amount,
		kindType:         kind,
		nonce:            nonce,
		signature:        signature,
	}
}

// SetUser sets the signing account's wallet address. Pass either this or
// SetSigner.
func (s *SubAccountTransferService) SetUser(user string) *SubAccountTransferService {
	s.user = user
	return s
}

// SetSigner sets the approved agent wallet authorizing the transfer on the
// user's behalf. Pass either this or SetUser.
func (s *SubAccountTransferService) SetSigner(signer string) *SubAccountTransferService {
	s.signer = signer
	return s
}

// SetFromAccountAddress sets the source wallet address. Required when the
// source account differs from the signing account.
func (s *SubAccountTransferService) SetFromAccountAddress(addr string) *SubAccountTransferService {
	s.fromAccountAddress = addr
	return s
}

func (s *SubAccountTransferService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	params := map[string]string{
		"toAccountAddress": s.toAccountAddress,
		"asset":            s.asset,
		"amount":           s.amount,
		"kindType":         string(s.kindType),
		"nonce":            formatInt64(s.nonce),
		"signature":        s.signature,
	}
	if s.user != "" {
		params["user"] = s.user
	}
	if s.signer != "" {
		params["signer"] = s.signer
	}
	if s.fromAccountAddress != "" {
		params["fromAccountAddress"] = s.fromAccountAddress
	}
	req := request.Post(s.c, ctx, "/fapi/v3/subAccountTransfer", params)
	return request.Do[GenericCodeMsg](req)
}

// MigrateUserAssetsService -- POST /fapi/v3/asset/migrateUser (WITHDRAW)
//
// Migrates every positive-balance asset from a source account into the
// authenticated user's account. The signature must be produced with the
// *source* user's wallet private key (not the configured signer), so -- like
// the sub-account flows -- the caller supplies the already-computed signature.
// Message body: user={user}&nonce={nonce}
//
// The source account must have no open positions and no open orders, and up to
// 300 assets are processed per batch. Record the returned batchId to query
// status via MigrateUserAssetsHistoryService.
type MigrateUserAssetsService struct {
	c         *FuturesClient
	user      string
	nonce     int64
	signature string
}

func (c *FuturesClient) NewMigrateUserAssetsService(user string, nonce int64, signature string) *MigrateUserAssetsService {
	return &MigrateUserAssetsService{c: c, user: user, nonce: nonce, signature: signature}
}

func (s *MigrateUserAssetsService) Do(ctx context.Context) (*MigrateUserAssetsResponse, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/asset/migrateUser", map[string]string{
		"user":      s.user,
		"nonce":     formatInt64(s.nonce),
		"signature": s.signature,
	})
	return request.Do[MigrateUserAssetsResponse](req)
}

// MigrateUserAssetsResponse carries the batchId used to query migration status.
// batchId is empty when the source account had no positive-balance assets.
type MigrateUserAssetsResponse struct {
	BatchID string `json:"batchId"`
}

// MigrateUserAssetsHistoryService -- GET /fapi/v3/asset/migrateUser/history (USER_DATA)
//
// Queries a migration batch by batchId. Signs with the regular configured
// signer, so this calls WithSignature.
type MigrateUserAssetsHistoryService struct {
	c       *FuturesClient
	batchID string
}

func (c *FuturesClient) NewMigrateUserAssetsHistoryService(batchID string) *MigrateUserAssetsHistoryService {
	return &MigrateUserAssetsHistoryService{c: c, batchID: batchID}
}

func (s *MigrateUserAssetsHistoryService) Do(ctx context.Context) (*MigrateUserAssetsHistory, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/asset/migrateUser/history", map[string]string{
		"batchId": s.batchID,
	}).WithSignature()
	return request.Do[MigrateUserAssetsHistory](req)
}

type MigrateUserAssetsHistory struct {
	BatchID         string                    `json:"batchId"`
	TotalCount      int                       `json:"totalCount"`
	SuccessCount    int                       `json:"successCount"`
	ProcessingCount int                       `json:"processingCount"`
	FailCount       int                       `json:"failCount"`
	InitCount       int                       `json:"initCount"`
	Details         []MigrateUserAssetsDetail `json:"details"`
}

// MigrateUserAssetsDetail is one per-asset migration record. status uses
// I=pending, S=success, F=failed; fromStatus/toStatus use S/F/P (processing).
// tranId is null until the record is processed.
type MigrateUserAssetsDetail struct {
	ID            int64           `json:"id"`
	Asset         string          `json:"asset"`
	Amount        decimal.Decimal `json:"amount"`
	TranID        *int64          `json:"tranId"`
	Status        string          `json:"status"`
	FromStatus    *string         `json:"fromStatus"`
	FromErrorCode *string         `json:"fromErrorCode"`
	FromResponse  *string         `json:"fromResponse"`
	ToStatus      *string         `json:"toStatus"`
	ToErrorCode   *string         `json:"toErrorCode"`
	ToResponse    *string         `json:"toResponse"`
}

func formatInt64(i int64) string {
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	n := len(buf)
	for i > 0 {
		n--
		buf[n] = digits[i%10]
		i /= 10
	}
	if neg {
		n--
		buf[n] = '-'
	}
	return string(buf[n:])
}
