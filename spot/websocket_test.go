package spot

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
)

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

// liveDepth5Frame is a captured sstream btcusdt@depth5 frame. The docs show
// partial depth levels as "bids"/"asks"; the stream sends "b"/"a".
const liveDepth5Frame = `{"e":"depthUpdate","E":1790603482674,"T":1790603482622,"s":"BTCUSDT",` +
	`"U":6650314299,"u":6650314358,"pu":6650314201,"b":[["83549.21","1.75275"],["83548.38","0.00038"],` +
	`["83548.37","0.01913"],["83548.36","0.00843"],["83548.28","0.00121"]],"a":[["83550.04","0.70332"],` +
	`["83552.01","0.00344"],["83554.01","0.00207"],["83555.53","0.00033"],["83555.98","0.00237"]]}`

func TestPartialDepthEventLevels(t *testing.T) {
	var ev WsPartialDepthEvent
	if err := common.JSONUnmarshal([]byte(liveDepth5Frame), &ev); err != nil {
		t.Fatal(err)
	}
	if len(ev.Bids) != 5 || len(ev.Asks) != 5 || ev.Bids[0] != [2]string{"83549.21", "1.75275"} ||
		ev.Asks[4] != [2]string{"83555.98", "0.00237"} || !ev.TransactionTime.Equal(time.UnixMilli(1790603482622)) {
		t.Errorf("decoded %+v", ev)
	}
}

// liveBookTickerFrame is a captured sstream btcusdt@bookTicker frame; the docs
// leave out e, E and T.
const liveBookTickerFrame = `{"u":6652628411,"e":"bookTicker","s":"BTCUSDT","b":"83335.59",` +
	`"B":"0.08198","a":"83336.42","A":"1.72705","T":1790610987719,"E":1790610987722}`

func TestBookTickerEventTimes(t *testing.T) {
	var ev WsBookTickerEvent
	if err := common.JSONUnmarshal([]byte(liveBookTickerFrame), &ev); err != nil {
		t.Fatal(err)
	}
	if ev.EventType != "bookTicker" || !ev.EventTime.Equal(time.UnixMilli(1790610987722)) ||
		!ev.TransactionTime.Equal(time.UnixMilli(1790610987719)) || ev.UpdateID != 6652628411 || ev.BidPrice.String() != "83335.59" {
		t.Errorf("decoded %+v", ev)
	}
}
