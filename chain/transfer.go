package chain

import (
	"context"
	"strconv"

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
