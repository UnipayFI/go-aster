package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	asterCommon "github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-json-experiment/json"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
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

// WithProxy routes all REST traffic through the given proxy. Supported
// schemes: http, https, socks5, socks5h. Pass userinfo in the URL for
// authenticated proxies, e.g. "http://user:pass@host:8080" or
// "socks5://user:pass@host:1080". Invalid URLs are logged and skipped.
func WithProxy(proxyURL string) Options {
	return func(opt *Option) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.client.SetProxy(proxyURL)
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithProxy: socks5 setup failed: %v", err)
				return
			}
			transport := cloneDefaultTransport()
			transport.Proxy = nil
			transport.DialContext = dialCtx
			opt.client.SetTransport(transport)
		default:
			opt.logger.Errorf("WithProxy: unsupported scheme %q (want http, https, socks5, socks5h)", u.Scheme)
		}
	}
}

type WebSocketOption struct {
	logger log.Logger
	client *resty.Client
	dialer *websocket.Dialer
}

type WebSocketOptions func(*WebSocketOption)

func defaultWebSocketOption() *WebSocketOption {
	return &WebSocketOption{
		logger: log.GetDefaultLogger(),
		client: defaultHttpClient(),
		dialer: defaultDialer(),
	}
}

func defaultDialer() *websocket.Dialer {
	return &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
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

// WithWebSocketProxy routes the WebSocket dial through the given proxy.
// Supported schemes: http, https, socks5, socks5h. Behaves like WithProxy:
// userinfo in the URL is passed through for auth, invalid URLs are logged
// and skipped.
func WithWebSocketProxy(proxyURL string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithWebSocketProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.dialer.Proxy = http.ProxyURL(u)
			opt.dialer.NetDialContext = nil
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithWebSocketProxy: socks5 setup failed: %v", err)
				return
			}
			opt.dialer.Proxy = nil
			opt.dialer.NetDialContext = dialCtx
		default:
			opt.logger.Errorf("WithWebSocketProxy: unsupported scheme %q (want http, https, socks5, socks5h)", u.Scheme)
		}
	}
}

// socks5DialContext builds a DialContext that tunnels TCP through the SOCKS5
// proxy described by u. socks5h is accepted as an alias of socks5: the SOCKS5
// dialer in golang.org/x/net/proxy already resolves hostnames remotely.
func socks5DialContext(u *url.URL) (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	su := *u
	if strings.EqualFold(su.Scheme, "socks5h") {
		su.Scheme = "socks5"
	}
	pd, err := proxy.FromURL(&su, proxy.Direct)
	if err != nil {
		return nil, err
	}
	if cd, ok := pd.(proxy.ContextDialer); ok {
		return cd.DialContext, nil
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return pd.Dial(network, addr)
	}, nil
}

func cloneDefaultTransport() *http.Transport {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		return t.Clone()
	}
	// Should never happen; guard anyway so SetTransport always gets a usable type.
	panic(fmt.Sprintf("aster: http.DefaultTransport is not *http.Transport (got %T)", http.DefaultTransport))
}
