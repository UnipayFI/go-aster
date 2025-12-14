package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/client"
)

type SpotClient struct {
	*client.Client
}

func NewSpotClient(options ...client.Options) *SpotClient {
	return &SpotClient{
		client.NewClient(options...),
	}
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
