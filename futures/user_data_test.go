package futures

import (
	"testing"
	"time"
)

// TestDecodeListenKeyExpired checks the expiry event with E as a JSON number
// (as documented) and as a quoted number (as Binance has sent it).
func TestDecodeListenKeyExpired(t *testing.T) {
	want := time.UnixMilli(1576653824250)
	for _, msg := range []string{
		`{"e":"listenKeyExpired","E":1576653824250}`,
		`{"e":"listenKeyExpired","E":"1576653824250"}`,
	} {
		ev, err := decodeUserDataEvent([]byte(msg))
		if err != nil {
			t.Fatalf("%s: %v", msg, err)
		}
		if ev.ListenKeyExpired == nil || !ev.ListenKeyExpired.EventTime.Equal(want) || !ev.EventTime.Equal(want) {
			t.Errorf("%s: decoded as %+v", msg, ev)
		}
	}
	if _, err := decodeUserDataEvent([]byte(`{"e":"listenKeyExpired","E":"soon"}`)); err == nil {
		t.Error("invalid E: want error")
	}
}
