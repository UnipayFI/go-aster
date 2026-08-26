// Package chain wraps the Aster-Chain REST API (/aster-chain/v3/*), which is
// served from its own host rather than the spot or futures gateways. It covers
// the on-chain side of an Aster account: deposit addresses, withdrawals to
// EVM/Solana addresses, wallet transfers between the spot and perp accounts,
// and withdrawal limit/fee queries.
package chain

import (
	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/request"
)

var _ request.Client = (*ChainClient)(nil)

type ChainClient struct {
	*client.Client
}

// NewChainClient constructs a REST client for /aster-chain/v3/* endpoints.
//
// Aster-Chain publishes no server-time endpoint, so unlike the spot and
// futures clients this one has no SyncServerTime. Nonces still come from the
// shared generator, so client.WithTimeOffset (or a SyncServerTime call on a
// spot/futures client sharing the same clock) is the way to correct drift.
func NewChainClient(options ...client.Options) *ChainClient {
	return &ChainClient{client.NewClient(client.ProductChain, options...)}
}
