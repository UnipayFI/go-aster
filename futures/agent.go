package futures

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-aster/v3/request"
)

// RegisterAndApproveAgentService -- POST /fapi/v3/registerAndApproveAgent (PUBLIC)
//
// Registers an API agent account and grants trading/withdrawal permissions in
// a single atomic call. This endpoint is unauthenticated -- no API key or HMAC
// header is required -- and all authorization is verified through the supplied
// signature, so (like the sub-account flows) the caller computes the signature
// and passes it in. The SDK does NOT call WithSignature here.
//
// The signature is produced by signing this canonical message body with the
// user's wallet private key:
//
//	user={user}&nonce={nonce}&agentName={agentName}&agentAddress={agentAddress}&expired={expired}&signatureChainId={signatureChainId}&canSpotTrade={canSpotTrade}&canPerpTrade={canPerpTrade}&canWithdraw={canWithdraw}&ipWhitelist={ipWhitelist}
//
// EVM addresses wrap that string as message.msg in the documented EIP-712
// typed data (signatureChainId=56); Solana addresses sign with Ed25519
// (signatureChainId=101). ipWhitelist is required and must not be empty when
// canWithdraw is true.
type RegisterAndApproveAgentService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewRegisterAndApproveAgentService(user, agentName, agentAddress string, nonce, expired, signatureChainId int64, canSpotTrade, canPerpTrade, canWithdraw bool, signature string) *RegisterAndApproveAgentService {
	return &RegisterAndApproveAgentService{c: c, params: map[string]string{
		"user":             user,
		"nonce":            strconv.FormatInt(nonce, 10),
		"agentName":        agentName,
		"agentAddress":     agentAddress,
		"expired":          strconv.FormatInt(expired, 10),
		"signatureChainId": strconv.FormatInt(signatureChainId, 10),
		"canSpotTrade":     strconv.FormatBool(canSpotTrade),
		"canPerpTrade":     strconv.FormatBool(canPerpTrade),
		"canWithdraw":      strconv.FormatBool(canWithdraw),
		"signature":        signature,
	}}
}

// SetIPWhitelist sets the space-separated list of permitted IPs/CIDR ranges.
// Required and must not be empty when canWithdraw is true.
func (s *RegisterAndApproveAgentService) SetIPWhitelist(ipWhitelist string) *RegisterAndApproveAgentService {
	s.params["ipWhitelist"] = ipWhitelist
	return s
}

// SetAgentCode sets an optional referral/invitation code.
func (s *RegisterAndApproveAgentService) SetAgentCode(code string) *RegisterAndApproveAgentService {
	s.params["agentCode"] = code
	return s
}

func (s *RegisterAndApproveAgentService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/registerAndApproveAgent", s.params)
	return request.Do[GenericCodeMsg](req)
}
