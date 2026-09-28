package futures

import "testing"

// TestStreamPath locks the stream names Aster accepts: the symbol is
// lowercased, event suffixes keep their case, and all-market streams are
// passed through untouched (a lowercased "!markprice@arr" gets no data).
func TestStreamPath(t *testing.T) {
	tests := map[string]string{
		"BTCUSDT@aggTrade":      "/ws/btcusdt@aggTrade",
		"BTCUSDT@markPrice@1s":  "/ws/btcusdt@markPrice@1s",
		"BTCUSDT@kline_1m":      "/ws/btcusdt@kline_1m",
		"BTCUSDT@bookTicker":    "/ws/btcusdt@bookTicker",
		"BTCUSDT@depth5@100ms":  "/ws/btcusdt@depth5@100ms",
		"!markPrice@arr":        "/ws/!markPrice@arr",
		"!markPrice@arr@1s":     "/ws/!markPrice@arr@1s",
		"!miniTicker@arr":       "/ws/!miniTicker@arr",
		"!ticker@arr":           "/ws/!ticker@arr",
		"!bookTicker":           "/ws/!bookTicker",
		"!forceOrder@arr":       "/ws/!forceOrder@arr",
		"ETHUSDT@forceOrder":    "/ws/ethusdt@forceOrder",
		"ETHUSDT@miniTicker":    "/ws/ethusdt@miniTicker",
		"ETHUSDT@ticker":        "/ws/ethusdt@ticker",
		"ETHUSDT@depth@500ms":   "/ws/ethusdt@depth@500ms",
		"ETHUSDT@depth20@500ms": "/ws/ethusdt@depth20@500ms",
		"ethusdt@aggTrade":      "/ws/ethusdt@aggTrade",
	}
	for in, want := range tests {
		if got := streamPath(in); got != want {
			t.Errorf("streamPath(%q) = %q, want %q", in, got, want)
		}
	}
}
