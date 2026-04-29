package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/request"
)

var _ request.Client = (*FuturesClient)(nil)

type FuturesClient struct {
	*client.Client
}

func NewFuturesClient(options ...client.Options) *FuturesClient {
	options = append(
		[]client.Options{client.WithBaseURL(common.DEFAULT_FUTURES_BASE_URL)},
		options...,
	)
	return &FuturesClient{client.NewClient(options...)}
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
	options = append(
		[]client.WebSocketOptions{client.WithWebSocketBaseURL(common.DEFAULT_FUTURES_WEBSOCKET_BASE_URL)},
		options...,
	)
	return &FuturesWebSocketClient{client.NewWebSocketClient(options...)}
}
