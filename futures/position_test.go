package futures

import (
	"context"
	"testing"
	"time"
)

// TestAssetExchangeResponse checks the documented response and the empty
// body Aster returns when no exchange took place.
func TestAssetExchangeResponse(t *testing.T) {
	got, err := newTestClient(t, `{"accountId":123456,"asset":"USDT","balance":"100.50","updateTime":1774077572319}`).
		NewAssetExchangeService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.AccountID != 123456 || got.Asset != "USDT" || got.Balance.String() != "100.5" ||
		!got.UpdateTime.Equal(time.UnixMilli(1774077572319)) {
		t.Errorf("decoded %+v", got)
	}

	for _, body := range []string{"", " \n"} {
		got, err := newTestClient(t, body).NewAssetExchangeService().Do(context.Background())
		if err != nil || got == nil || *got != (AssetExchangeResponse{}) {
			t.Errorf("body %q: got %+v, %v, want the zero response", body, got, err)
		}
	}
	if _, err := newTestClient(t, `{"updateTime":"x"}`).NewAssetExchangeService().Do(context.Background()); err == nil {
		t.Error("bad updateTime: want error")
	}
}
