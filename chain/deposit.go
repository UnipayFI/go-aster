package chain

import (
	"context"

	"github.com/UnipayFI/go-aster/v3/request"
)

// GetSpotDepositAddressService -- GET /aster-chain/v3/spot/user-deposit-address (USER_DATA)
//
// Returns the caller's dedicated deposit address on the given network. Only
// the spot account is supported, which is why SUI is the default network: SUI
// deposits are plain transfers to this address and need no contract call,
// whereas EVM chains deposit through the depositFor contract method and Solana
// through the depositSol / depositToken program methods.
type GetSpotDepositAddressService struct {
	c      *ChainClient
	params map[string]string
}

func (c *ChainClient) NewGetSpotDepositAddressService() *GetSpotDepositAddressService {
	return &GetSpotDepositAddressService{c: c, params: map[string]string{}}
}

// SetNetwork overrides the network whose deposit address is returned.
// Defaults to SUI.
func (s *GetSpotDepositAddressService) SetNetwork(network string) *GetSpotDepositAddressService {
	s.params["network"] = network
	return s
}

func (s *GetSpotDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/aster-chain/v3/spot/user-deposit-address", s.params).WithSignature()
	return request.Do[DepositAddress](req)
}

type DepositAddress struct {
	Network string `json:"network"`
	Address string `json:"address"`
}
