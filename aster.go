// Package aster is the entry point of the Aster DEX V3 Go SDK.
//
// Install: go get -u github.com/UnipayFI/go-aster/v3
// Import:  import aster "github.com/UnipayFI/go-aster/v3"
//
// V3 authenticates with EIP-712 + ECDSA (the API wallet model) instead of
// the V1 HMAC-SHA256 scheme. Configure credentials with
// client.WithAuth(userAddress, signerPrivateKeyHex); the signer address is
// derived from the private key.
package aster

import (
	"github.com/UnipayFI/go-aster/v3/chain"
	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/futures"
	"github.com/UnipayFI/go-aster/v3/spot"
)

// NewSpotClient constructs a REST client for /api/v3/* spot endpoints.
func NewSpotClient(options ...client.Options) *spot.SpotClient {
	return spot.NewSpotClient(options...)
}

// NewSpotWebSocketClient constructs a WebSocket client for spot market and
// user-data streams.
func NewSpotWebSocketClient(options ...client.WebSocketOptions) *spot.SpotWebSocketClient {
	return spot.NewSpotWebSocketClient(options...)
}

// NewFuturesClient constructs a REST client for /fapi/v3/* futures endpoints.
func NewFuturesClient(options ...client.Options) *futures.FuturesClient {
	return futures.NewFuturesClient(options...)
}

// NewFuturesWebSocketClient constructs a WebSocket client for futures market
// and user-data streams.
func NewFuturesWebSocketClient(options ...client.WebSocketOptions) *futures.FuturesWebSocketClient {
	return futures.NewFuturesWebSocketClient(options...)
}

// NewChainClient constructs a REST client for /aster-chain/v3/* endpoints
// (deposit addresses, withdrawals, and spot/perp wallet transfers). These are
// served from their own host, so this client is separate from the spot and
// futures ones.
func NewChainClient(options ...client.Options) *chain.ChainClient {
	return chain.NewChainClient(options...)
}
