package futures

import (
	"github.com/UnipayFI/go-aster/client"
	"github.com/UnipayFI/go-aster/internal/request"
)

var _ request.Client = (*FuturesClient)(nil)

type FuturesClient struct {
	*client.Client
}

func NewFuturesClient(options ...client.Options) *FuturesClient {
	return &FuturesClient{
		client.NewClient(options...),
	}
}
