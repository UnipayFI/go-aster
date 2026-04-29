package common

import "time"

const (
	GO_ASTER_USER_AGENT = "go-aster/3.0"

	DEFAULT_SPOT_BASE_URL              = "https://sapi.asterdex.com"
	DEFAULT_FUTURES_BASE_URL           = "https://fapi.asterdex.com"
	DEFAULT_SPOT_WEBSOCKET_BASE_URL    = "wss://sstream.asterdex.com"
	DEFAULT_FUTURES_WEBSOCKET_BASE_URL = "wss://fstream.asterdex.com"

	WEBSOCKET_STREAM_SEPARATOR   = "/ws/"
	WEBSOCKET_COMBINED_SEPARATOR = "/stream?streams="

	DEFAULT_KEEP_ALIVE_INTERVAL = 10 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second

	// V3 EIP712 signing
	EIP712_DOMAIN_NAME         = "AsterSignTransaction"
	EIP712_DOMAIN_VERSION      = "1"
	EIP712_VERIFYING_CONTRACT  = "0x0000000000000000000000000000000000000000"
	EIP712_PRIMARY_TYPE        = "Message"
	EIP712_CHAIN_ID_MAINNET    = int64(1666)
	EIP712_CHAIN_ID_TESTNET    = int64(714)
	DEFAULT_EIP712_CHAIN_ID    = EIP712_CHAIN_ID_MAINNET
)
