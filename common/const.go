package common

import "time"

const (
	GO_ASTER_USER_AGENT = "go-aster/3.0"

	// Mainnet endpoints.
	DEFAULT_SPOT_BASE_URL_MAINNET              = "https://sapi.asterdex.com"
	DEFAULT_FUTURES_BASE_URL_MAINNET           = "https://fapi.asterdex.com"
	DEFAULT_SPOT_WEBSOCKET_BASE_URL_MAINNET    = "wss://sstream.asterdex.com"
	DEFAULT_FUTURES_WEBSOCKET_BASE_URL_MAINNET = "wss://fstream.asterdex.com"
	DEFAULT_CHAIN_BASE_URL                     = "https://chainapi.asterdex.com"

	// Testnet endpoints.
	DEFAULT_SPOT_BASE_URL_TESTNET              = "https://sapi.asterdex-testnet.com"
	DEFAULT_FUTURES_BASE_URL_TESTNET           = "https://fapi.asterdex-testnet.com"
	DEFAULT_SPOT_WEBSOCKET_BASE_URL_TESTNET    = "wss://sstream.asterdex-testnet.com"
	DEFAULT_FUTURES_WEBSOCKET_BASE_URL_TESTNET = "wss://fstream.asterdex-testnet.com"

	WEBSOCKET_STREAM_SEPARATOR   = "/ws/"
	WEBSOCKET_COMBINED_SEPARATOR = "/stream?streams="

	DEFAULT_KEEP_ALIVE_INTERVAL = 10 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second

	// V3 EIP712 signing
	EIP712_DOMAIN_NAME        = "AsterSignTransaction"
	EIP712_DOMAIN_VERSION     = "1"
	EIP712_VERIFYING_CONTRACT = "0x0000000000000000000000000000000000000000"
	EIP712_PRIMARY_TYPE       = "Message"
	EIP712_CHAIN_ID_MAINNET   = int64(1666)
	EIP712_CHAIN_ID_TESTNET   = int64(714)
)

// Network identifies which Aster network a client targets. Both the REST/WS
// base URLs and the EIP-712 chainId are derived from this value.
type Network int

const (
	Mainnet Network = iota
	Testnet
)

// ChainID returns the EIP-712 chainId for this network.
func (n Network) ChainID() int64 {
	if n == Testnet {
		return EIP712_CHAIN_ID_TESTNET
	}
	return EIP712_CHAIN_ID_MAINNET
}

// SpotBaseURL returns the spot REST base URL for this network.
func (n Network) SpotBaseURL() string {
	if n == Testnet {
		return DEFAULT_SPOT_BASE_URL_TESTNET
	}
	return DEFAULT_SPOT_BASE_URL_MAINNET
}

// FuturesBaseURL returns the futures REST base URL for this network.
func (n Network) FuturesBaseURL() string {
	if n == Testnet {
		return DEFAULT_FUTURES_BASE_URL_TESTNET
	}
	return DEFAULT_FUTURES_BASE_URL_MAINNET
}

// SpotWebSocketBaseURL returns the spot WebSocket base URL for this network.
func (n Network) SpotWebSocketBaseURL() string {
	if n == Testnet {
		return DEFAULT_SPOT_WEBSOCKET_BASE_URL_TESTNET
	}
	return DEFAULT_SPOT_WEBSOCKET_BASE_URL_MAINNET
}

// ChainBaseURL returns the Aster-Chain REST base URL. The docs publish a
// single host with no testnet counterpart, so both networks resolve to it;
// override with client.WithBaseURL if that changes.
func (Network) ChainBaseURL() string {
	return DEFAULT_CHAIN_BASE_URL
}

// FuturesWebSocketBaseURL returns the futures WebSocket base URL for this network.
func (n Network) FuturesWebSocketBaseURL() string {
	if n == Testnet {
		return DEFAULT_FUTURES_WEBSOCKET_BASE_URL_TESTNET
	}
	return DEFAULT_FUTURES_WEBSOCKET_BASE_URL_MAINNET
}
