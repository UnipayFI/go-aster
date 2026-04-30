package client

import (
	asterCommon "github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-json-experiment/json"
	"github.com/go-resty/resty/v2"
)

const (
	DEFAULT_RECV_WINDOW = 5000
)

type Option struct {
	userAddress         string
	signerAddress       string
	signerPrivateKeyHex string
	chainID             int64
	recvWindow          int64
	logger              log.Logger
	signFn              SignFn
	client              *resty.Client
	timeOffsetMs        int64
}

type Options func(*Option)

func defaultOption() *Option {
	return &Option{
		chainID:    asterCommon.DEFAULT_EIP712_CHAIN_ID,
		recvWindow: DEFAULT_RECV_WINDOW,
		logger:     log.GetDefaultLogger(),
		client:     defaultHttpClient(),
	}
}

func defaultHttpClient() *resty.Client {
	return resty.New().
		SetJSONMarshaler(func(v any) ([]byte, error) {
			return json.Marshal(v)
		}).
		SetJSONUnmarshaler(func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		})
}

func WithBaseURL(baseURL string) Options {
	return func(opt *Option) {
		opt.client.SetBaseURL(baseURL)
	}
}

func WithLogger(logger log.Logger) Options {
	return func(opt *Option) {
		opt.logger = logger
	}
}

// WithAuth configures the credentials used for V3 signed requests.
// userAddress is the main wallet address (may be empty when not required by
// the endpoint). signerPrivateKeyHex is the API wallet private key in hex
// (with or without the "0x" prefix); the signer address is derived from it.
func WithAuth(userAddress, signerPrivateKeyHex string) Options {
	return func(opt *Option) {
		opt.userAddress = userAddress
		opt.signerPrivateKeyHex = signerPrivateKeyHex
	}
}

// WithTEEAuth configures credentials for TEE / HSM / remote-signer mode where
// the signer's private key never leaves the enclave. Pass the main wallet
// address and the API-wallet (signer) address; no local private key is
// required. Must be combined with WithSignRequestFn, which is what actually
// performs the EIP-712 signing -- typically by shelling out to the TEE binary.
func WithTEEAuth(userAddress, signerAddress string) Options {
	return func(opt *Option) {
		opt.userAddress = userAddress
		opt.signerAddress = signerAddress
	}
}

// WithChainID overrides the EIP-712 domain chainId. Defaults to mainnet (1666).
// Use 714 for testnet.
func WithChainID(chainID int64) Options {
	return func(opt *Option) {
		opt.chainID = chainID
	}
}

// WithSignRequestFn replaces the default EIP-712 + ECDSA signer with a custom
// implementation (e.g. for HSM, TEE, or a remote signing service).
func WithSignRequestFn(signFn SignFn) Options {
	return func(opt *Option) {
		opt.signFn = signFn
	}
}

func WithRecvWindow(recvWindow int64) Options {
	return func(opt *Option) {
		opt.recvWindow = recvWindow
	}
}

func WithTimeOffset(timeOffsetMs int64) Options {
	return func(opt *Option) {
		opt.timeOffsetMs = timeOffsetMs
	}
}

type WebSocketOption struct {
	logger log.Logger
	client *resty.Client
}

type WebSocketOptions func(*WebSocketOption)

func defaultWebSocketOption() *WebSocketOption {
	return &WebSocketOption{
		logger: log.GetDefaultLogger(),
		client: defaultHttpClient(),
	}
}

func WithWebSocketBaseURL(baseURL string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		opt.client.SetBaseURL(baseURL)
	}
}

func WithWebSocketLogger(logger log.Logger) WebSocketOptions {
	return func(opt *WebSocketOption) {
		opt.logger = logger
	}
}
