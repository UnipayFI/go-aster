package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/client"
	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/request"
)

var _ request.Client = (*SpotClient)(nil)

type SpotClient struct {
	*client.Client
}

func NewSpotClient(options ...client.Options) *SpotClient {
	options = append(
		[]client.Options{client.WithBaseURL(common.DEFAULT_SPOT_BASE_URL)},
		options...,
	)
	return &SpotClient{client.NewClient(options...)}
}

// SyncServerTime synchronizes with server time and calculates offset
func (c *SpotClient) SyncServerTime(ctx context.Context) error {
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

type SpotWebSocketClient struct {
	*client.WebSocketClient
}

func NewSpotWebSocketClient(options ...client.WebSocketOptions) *SpotWebSocketClient {
	options = append(
		[]client.WebSocketOptions{client.WithWebSocketBaseURL(common.DEFAULT_SPOT_WEBSOCKET_BASE_URL)},
		options...,
	)
	return &SpotWebSocketClient{client.NewWebSocketClient(options...)}
}
