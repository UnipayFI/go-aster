package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/client"
	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/internal/request"
)

var _ request.Client = (*FuturesClient)(nil)

type FuturesClient struct {
	*client.Client
}

func NewFuturesClient(options ...client.Options) *FuturesClient {
	opts := append([]client.Options{client.WithBaseURL(common.DEFAULT_FUTURES_BASE_URL)}, options...)
	return &FuturesClient{
		client.NewClient(opts...),
	}
}

func (c *FuturesClient) SyncServerTime(ctx context.Context) error {
	localTimeBefore := time.Now().UnixMilli()

	serverTime, err := c.NewGetServerTimeService().Do(ctx)
	if err != nil {
		return err
	}

	localTimeAfter := time.Now().UnixMilli()
	localTime := (localTimeBefore + localTimeAfter) / 2

	c.SetTimeOffset(localTime - serverTime.ServerTime.UnixMilli())

	c.GetLogger().Infof("Time sync: local=%d, server=%d, offset=%dms",
		localTime, serverTime.ServerTime.UnixMilli(), c.GetTimeOffsetMs())

	return nil
}

type FuturesWebSocketClient struct {
	*client.WebSocketClient
}

func NewFuturesWebSocketClient(options ...client.Options) *FuturesWebSocketClient {
	opts := append([]client.Options{client.WithBaseURL(common.DEFAULT_FUTURES_WEBSOCKET_BASE_URL)}, options...)
	return &FuturesWebSocketClient{client.NewWebSocketClient(opts...)}
}
