package spot

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
)

// depthBody is a trimmed live /api/v3/depth response.
const depthBody = `{"lastUpdateId":6651296101,"E":1790606631282,"T":1790606631279,` +
	`"symbol":"ASTERUSDT","bids":[["0.69370","21.33"]],"asks":[["0.69400","10.00"]]}`

// TestGetDepth checks a normal book, the empty array Aster sends for symbols
// it no longer lists, and an API error.
func TestGetDepth(t *testing.T) {
	serve := func(status int, body string) *SpotClient {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}))
		t.Cleanup(srv.Close)
		return NewSpotClient(client.WithBaseURL(srv.URL))
	}

	got, err := serve(http.StatusOK, depthBody).NewGetDepthService("ASTERUSDT").Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.LastUpdateId != 6651296101 || !got.TransactionTime.Equal(time.UnixMilli(1790606631279)) || len(got.Bids) != 1 || len(got.Asks) != 1 {
		t.Errorf("depth = %+v", got)
	}

	if _, err := serve(http.StatusOK, "[]").NewGetDepthService("B2USDT").Do(context.Background()); !errors.Is(err, ErrNoOrderBook) {
		t.Errorf("empty array: err = %v, want ErrNoOrderBook", err)
	}

	_, err = serve(http.StatusBadRequest, `{"code":-1121,"msg":"Invalid symbol."}`).NewGetDepthService("ZZZ").Do(context.Background())
	if apiErr, ok := err.(*client.APIError); !ok || apiErr.Code != -1121 {
		t.Errorf("invalid symbol: err = %v, want APIError -1121", err)
	}
}
