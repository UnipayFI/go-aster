package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/request"
)

var _ request.Client = (*FuturesClient)(nil)

type FuturesClient struct {
	*client.Client
}

func NewFuturesClient(options ...client.Options) *FuturesClient {
	return &FuturesClient{client.NewClient(client.ProductFutures, options...)}
}

func (c *FuturesClient) SyncServerTime(ctx context.Context) error {
	localBefore := time.Now().UnixMilli()
	resp, err := c.NewGetServerTimeService().Do(ctx)
	if err != nil {
		return err
	}
	localAfter := time.Now().UnixMilli()
	local := (localBefore + localAfter) / 2
	c.SetTimeOffset(local - resp.ServerTime.UnixMilli())
	c.GetLogger().Infof("Time sync: local=%d, server=%d, offset=%dms",
		local, resp.ServerTime.UnixMilli(), c.GetTimeOffsetMs())
	return nil
}

type FuturesWebSocketClient struct {
	*client.WebSocketClient
}

func NewFuturesWebSocketClient(options ...client.WebSocketOptions) *FuturesWebSocketClient {
	return &FuturesWebSocketClient{client.NewWebSocketClient(client.ProductFutures, options...)}
}
