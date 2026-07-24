package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/request"
)

var _ request.Client = (*SpotClient)(nil)

type SpotClient struct {
	*client.Client
}

func NewSpotClient(options ...client.Options) *SpotClient {
	return &SpotClient{client.NewClient(client.ProductSpot, options...)}
}

// SyncServerTime aligns the client's microsecond nonce generator with the
// server clock. V3's nonce window is ±60 seconds, so a few hundred ms of
// drift is harmless, but on long-running processes it's still safer to
// resync periodically.
func (c *SpotClient) SyncServerTime(ctx context.Context) error {
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

type SpotWebSocketClient struct {
	*client.WebSocketClient
}

func NewSpotWebSocketClient(options ...client.WebSocketOptions) *SpotWebSocketClient {
	return &SpotWebSocketClient{client.NewWebSocketClient(client.ProductSpot, options...)}
}
