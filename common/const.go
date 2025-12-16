package common

import "time"

const (
	GO_ASTER_USER_AGENT                = "go-aster/1.0"
	DEFAULT_SPOT_BASE_URL              = "https://sapi.asterdex.com"
	DEFAULT_FUTURES_BASE_URL           = "https://fapi.asterdex.com"
	DEFAULT_SPOT_WEBSOCKET_BASE_URL    = "wss://sstream.asterdex.com/ws"
	DEFAULT_FUTURES_WEBSOCKET_BASE_URL = "wss://fstream.asterdex.com/ws"

	DEFAULT_KEEP_ALIVE_INTERVAL = 10 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second
)
