package spot

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-aster/v3/request"
	"github.com/shopspring/decimal"
)

// Spot Builder management endpoints (approveBuilder / updateBuilder /
// delBuilder) are signed with the user's *main wallet* key rather than the
// signer/agent key configured on the client, and each uses its own EIP-712
// primaryType instead of the usual Message wrapper. Main-wallet keys normally
// live outside an SDK process, so -- as with the futures sub-account flows --
// these services accept an already-computed signature. The caller is expected
// to:
//
//  1. Build the EIP-712 typed data for the operation. The domain is the usual
//     {name: "AsterSignTransaction", version: "1", verifyingContract: 0x0...0}
//     but domain.chainId is the signatureChainId request parameter, NOT the
//     network-wide 1666/714 chainId used elsewhere in V3.
//  2. Sign it with the user's main wallet private key. Typed-data field names
//     are Title Case even though the request parameters are camelCase:
//
//     ApproveBuilder(string Builder,string MaxFeeRate,string BuilderName,string AsterChain,string User,uint256 Nonce)
//     UpdateBuilder(string Builder,string MaxFeeRate,string AsterChain,string User,uint256 Nonce)
//     DelBuilder(string Builder,string AsterChain,string User,uint256 Nonce)
//
//  3. Pass the resulting hex signature to the service.
//
// These services therefore do NOT call WithSignature: no nonce/signer is
// injected and the payload arrives already signed. Fee rates are taken as
// strings so the bytes sent match the bytes the caller signed.
//
// asterChain is the Aster environment identifier. The docs type it as a plain
// STRING without enumerating accepted values, so it is passed through as-is.
//
// All three management endpoints answer with HTTP 200 and a zero-byte body,
// so Do reports only an error.

// ApproveBuilderService -- POST /api/v3/approveBuilder (SPOT_TRADE)
//
// Authorizes a Spot Builder to place orders on the user's behalf, capped at
// maxFeeRate. The equivalent can also be done in one shot while approving the
// agent, via the spotBuilder / maxSpotFeeRate / spotBuilderName parameters of
// futures ApproveAgentService.
type ApproveBuilderService struct {
	c                *SpotClient
	builder          string
	maxFeeRate       string
	builderName      string
	asterChain       string
	user             string
	nonce            int64
	signatureChainId int64
	signature        string
}

func (c *SpotClient) NewApproveBuilderService(builder, maxFeeRate, asterChain, user string, nonce, signatureChainId int64, signature string) *ApproveBuilderService {
	return &ApproveBuilderService{
		c:                c,
		builder:          builder,
		maxFeeRate:       maxFeeRate,
		asterChain:       asterChain,
		user:             user,
		nonce:            nonce,
		signatureChainId: signatureChainId,
		signature:        signature,
	}
}

// SetBuilderName sets the optional Builder label. It is part of the signed
// ApproveBuilder typed data, so pass the same value that was signed (the empty
// string when the signature covered an empty BuilderName).
func (s *ApproveBuilderService) SetBuilderName(name string) *ApproveBuilderService {
	s.builderName = name
	return s
}

func (s *ApproveBuilderService) Do(ctx context.Context) error {
	params := map[string]string{
		"builder":          s.builder,
		"maxFeeRate":       s.maxFeeRate,
		"asterChain":       s.asterChain,
		"user":             s.user,
		"nonce":            strconv.FormatInt(s.nonce, 10),
		"signatureChainId": strconv.FormatInt(s.signatureChainId, 10),
		"signature":        s.signature,
	}
	if s.builderName != "" {
		params["builderName"] = s.builderName
	}
	req := request.Post(s.c, ctx, "/api/v3/approveBuilder", params)
	_, err := request.Do[struct{}](req)
	return err
}

// UpdateBuilderService -- POST /api/v3/updateBuilder (SPOT_TRADE)
//
// Changes the maxFeeRate cap of an already-approved Spot Builder. There is no
// builderName parameter here; the name set at approval time is kept. Returns
// error code -1130 ("No builder found") when the address was never approved.
type UpdateBuilderService struct {
	c                *SpotClient
	builder          string
	maxFeeRate       string
	asterChain       string
	user             string
	nonce            int64
	signatureChainId int64
	signature        string
}

func (c *SpotClient) NewUpdateBuilderService(builder, maxFeeRate, asterChain, user string, nonce, signatureChainId int64, signature string) *UpdateBuilderService {
	return &UpdateBuilderService{
		c:                c,
		builder:          builder,
		maxFeeRate:       maxFeeRate,
		asterChain:       asterChain,
		user:             user,
		nonce:            nonce,
		signatureChainId: signatureChainId,
		signature:        signature,
	}
}

func (s *UpdateBuilderService) Do(ctx context.Context) error {
	req := request.Post(s.c, ctx, "/api/v3/updateBuilder", map[string]string{
		"builder":          s.builder,
		"maxFeeRate":       s.maxFeeRate,
		"asterChain":       s.asterChain,
		"user":             s.user,
		"nonce":            strconv.FormatInt(s.nonce, 10),
		"signatureChainId": strconv.FormatInt(s.signatureChainId, 10),
		"signature":        s.signature,
	})
	_, err := request.Do[struct{}](req)
	return err
}

// DelBuilderService -- DELETE /api/v3/delBuilder (SPOT_TRADE)
//
// Revokes a Spot Builder's authorization. Returns error code -1130 ("No
// builder found") when the address was never approved. Futures Builders are
// revoked through a different endpoint (/fapi/v3/builder).
type DelBuilderService struct {
	c                *SpotClient
	builder          string
	asterChain       string
	user             string
	nonce            int64
	signatureChainId int64
	signature        string
}

func (c *SpotClient) NewDelBuilderService(builder, asterChain, user string, nonce, signatureChainId int64, signature string) *DelBuilderService {
	return &DelBuilderService{
		c:                c,
		builder:          builder,
		asterChain:       asterChain,
		user:             user,
		nonce:            nonce,
		signatureChainId: signatureChainId,
		signature:        signature,
	}
}

func (s *DelBuilderService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v3/delBuilder", map[string]string{
		"builder":          s.builder,
		"asterChain":       s.asterChain,
		"user":             s.user,
		"nonce":            strconv.FormatInt(s.nonce, 10),
		"signatureChainId": strconv.FormatInt(s.signatureChainId, 10),
		"signature":        s.signature,
	})
	_, err := request.Do[struct{}](req)
	return err
}

// GetBuildersService -- GET /api/v3/builder (USER_DATA)
//
// Lists the Spot Builders the authenticated user has approved. Unlike the
// management endpoints above this one uses the ordinary signer/agent EIP-712
// Message signature, so the SDK builds it via WithSignature.
type GetBuildersService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetBuildersService() *GetBuildersService {
	return &GetBuildersService{c: c}
}

func (s *GetBuildersService) Do(ctx context.Context) ([]SpotBuilder, error) {
	req := request.Get(ctx, s.c, "/api/v3/builder").WithSignature()
	resp, err := request.Do[[]SpotBuilder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotBuilder is one Spot Builder the user has approved. maxFeeRate comes back
// as a JSON number here (the approve/update requests take it as a string).
// builderName is "" when the approval carried no name.
type SpotBuilder struct {
	UserAddress    string          `json:"userAddress"`
	BuilderAddress string          `json:"builderAddress"`
	MaxFeeRate     decimal.Decimal `json:"maxFeeRate"`
	BuilderName    string          `json:"builderName"`
}
