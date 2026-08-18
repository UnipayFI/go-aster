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
// (signatureChainId=101). The EIP-712 domain.chainId MUST equal the
// signatureChainId request parameter (56 for EVM, 101 for Solana) -- it is no
// longer the network-wide 1666/714 chainId used by the other V3 flows.
// ipWhitelist is required and must not be empty when canWithdraw is true.
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

// ApproveAgentService -- POST /fapi/v3/approveAgent (USER_DATA)
//
// Grants an already-existing agent (API wallet) trading/withdrawal permissions
// on the user's account, and can authorize a Futures Builder and/or a Spot
// Builder in the same call. Use RegisterAndApproveAgentService instead when the
// agent still has to be registered.
//
// Like the sub-account flows, the signature is produced with the user's *main
// wallet* key rather than the configured signer key, so the caller supplies it
// and this service does NOT call WithSignature. The signature is EIP-712 with
// primaryType=ApproveAgent over the request parameters with their first letter
// upper-cased (agentName -> AgentName), per the official demo; the docs do not
// publish an explicit field list. Unlike registerAndApproveAgent, there is no
// signatureChainId parameter here -- domain.chainId is the network-wide 1666
// (mainnet) / 714 (testnet).
//
// Builder authorization is optional and independent per product: set builder +
// maxFeeRate for a Futures Builder, spotBuilder + maxSpotFeeRate for a Spot
// Builder. The Spot Builder can equivalently be authorized afterwards through
// spot ApproveBuilderService. Fee rates are taken as strings so the bytes sent
// match the bytes the caller signed.
type ApproveAgentService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewApproveAgentService(user, agentName, agentAddress string, nonce, expired int64, canSpotTrade, canPerpTrade, canWithdraw bool, signature string) *ApproveAgentService {
	return &ApproveAgentService{c: c, params: map[string]string{
		"user":         user,
		"nonce":        strconv.FormatInt(nonce, 10),
		"agentName":    agentName,
		"agentAddress": agentAddress,
		"expired":      strconv.FormatInt(expired, 10),
		"canSpotTrade": strconv.FormatBool(canSpotTrade),
		"canPerpTrade": strconv.FormatBool(canPerpTrade),
		"canWithdraw":  strconv.FormatBool(canWithdraw),
		"signature":    signature,
	}}
}

// SetIPWhitelist restricts the agent to the given space-separated IPs/CIDR
// ranges -- typically the Builder backend's egress IP.
func (s *ApproveAgentService) SetIPWhitelist(ipWhitelist string) *ApproveAgentService {
	s.params["ipWhitelist"] = ipWhitelist
	return s
}

// SetBuilder also authorizes this Futures Builder address. SetMaxFeeRate is
// required alongside it.
func (s *ApproveAgentService) SetBuilder(builder string) *ApproveAgentService {
	s.params["builder"] = builder
	return s
}

// SetMaxFeeRate caps the fee rate the Futures Builder may charge.
func (s *ApproveAgentService) SetMaxFeeRate(maxFeeRate string) *ApproveAgentService {
	s.params["maxFeeRate"] = maxFeeRate
	return s
}

// SetBuilderName labels the Futures Builder.
func (s *ApproveAgentService) SetBuilderName(name string) *ApproveAgentService {
	s.params["builderName"] = name
	return s
}

// SetSpotBuilder also authorizes this Spot Builder address. SetMaxSpotFeeRate
// is required alongside it.
func (s *ApproveAgentService) SetSpotBuilder(spotBuilder string) *ApproveAgentService {
	s.params["spotBuilder"] = spotBuilder
	return s
}

// SetMaxSpotFeeRate caps the fee rate the Spot Builder may charge.
func (s *ApproveAgentService) SetMaxSpotFeeRate(maxSpotFeeRate string) *ApproveAgentService {
	s.params["maxSpotFeeRate"] = maxSpotFeeRate
	return s
}

// SetSpotBuilderName labels the Spot Builder.
func (s *ApproveAgentService) SetSpotBuilderName(name string) *ApproveAgentService {
	s.params["spotBuilderName"] = name
	return s
}

func (s *ApproveAgentService) Do(ctx context.Context) (*GenericCodeMsg, error) {
	req := request.Post(s.c, ctx, "/fapi/v3/approveAgent", s.params)
	return request.Do[GenericCodeMsg](req)
}
