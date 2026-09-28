package spot

import "testing"

// TestStreamPath locks the stream names Aster accepts: the symbol is
// lowercased, event suffixes keep their case, and all-market streams are
// passed through untouched (a lowercased "!miniticker@arr" gets no data).
func TestStreamPath(t *testing.T) {
	tests := map[string]string{
		"BTCUSDT@aggTrade":     "/ws/btcusdt@aggTrade",
		"BTCUSDT@trade":        "/ws/btcusdt@trade",
		"BTCUSDT@kline_1m":     "/ws/btcusdt@kline_1m",
		"BTCUSDT@miniTicker":   "/ws/btcusdt@miniTicker",
		"BTCUSDT@ticker":       "/ws/btcusdt@ticker",
		"BTCUSDT@bookTicker":   "/ws/btcusdt@bookTicker",
		"BTCUSDT@depth5@100ms": "/ws/btcusdt@depth5@100ms",
		"BTCUSDT@depth@100ms":  "/ws/btcusdt@depth@100ms",
		"ASTERUSDT@tradepro":   "/ws/asterusdt@tradepro",
		"!miniTicker@arr":      "/ws/!miniTicker@arr",
		"!ticker@arr":          "/ws/!ticker@arr",
		"!bookTicker":          "/ws/!bookTicker",
		"asterusdt@aggTrade":   "/ws/asterusdt@aggTrade",
	}
	for in, want := range tests {
		if got := streamPath(in); got != want {
			t.Errorf("streamPath(%q) = %q, want %q", in, got, want)
		}
	}
}
