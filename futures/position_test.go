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

// docPositionRiskBody is the Hedge Mode response sample from the Position
// Information v3 doc, its missing commas restored.
const docPositionRiskBody = `[{"entryPrice":"6563.66500","marginType":"isolated","isAutoAddMargin":"false",` +
	`"isolatedMargin":"15517.54150468","leverage":"10","liquidationPrice":"5930.78","markPrice":"6679.50671178",` +
	`"maxNotionalValue":"20000000","positionAmt":"20.000","symbol":"BTCUSDT","unRealizedProfit":"2316.83423560",` +
	`"positionSide":"LONG","notional":"133590.13423560","isolatedWallet":"15517.54150468","updateTime":1625474304765},` +
	`{"entryPrice":"0.00000","marginType":"isolated","isAutoAddMargin":"false","isolatedMargin":"5413.95799991",` +
	`"leverage":"10","liquidationPrice":"7189.95","markPrice":"6679.50671178","maxNotionalValue":"20000000",` +
	`"positionAmt":"-10.000","symbol":"BTCUSDT","unRealizedProfit":"-1156.46711780","positionSide":"SHORT",` +
	`"notional":"-66795.0671178","isolatedWallet":"5413.95799991","updateTime":0}]`

func TestPositionRiskNotional(t *testing.T) {
	got, err := newTestClient(t, docPositionRiskBody).NewGetPositionRiskService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(got) != 2 || got[0].Notional.String() != "133590.1342356" || got[0].IsolatedWallet.String() != "15517.54150468" ||
		got[1].Notional.String() != "-66795.0671178" || !got[1].UpdateTime.Equal(time.UnixMilli(0)) {
		t.Errorf("decoded %+v", got)
	}
}

// docMarginHistoryBody is the response sample from the Get Position Margin
// Change History doc: only the first entry carries deltaType and clientTranId.
const docMarginHistoryBody = `[{"amount":"23.36332311","asset":"USDT","symbol":"BTCUSDT","time":1578047897183,` +
	`"type":1,"positionSide":"BOTH","deltaType":"TRADE","clientTranId":""},` +
	`{"amount":"100","asset":"USDT","symbol":"BTCUSDT","time":1578047900425,"type":1,"positionSide":"LONG"}]`

func TestPositionMarginHistoryDeltaType(t *testing.T) {
	got, err := newTestClient(t, docMarginHistoryBody).NewGetPositionMarginHistoryService("BTCUSDT").Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(got) != 2 || got[0].DeltaType != "TRADE" || got[0].ClientTranID != "" || got[1].DeltaType != "" ||
		!got[1].Time.Equal(time.UnixMilli(1578047900425)) {
		t.Errorf("decoded %+v", got)
	}
}
