package futures

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/UnipayFI/go-aster/v3/common"
)

// liveSymbol is BTCUSDT from a live /fapi/v3/exchangeInfo (filters dropped).
// The docs name the order types "OrderType"; the API sends "orderTypes".
const liveSymbol = `{"symbol":"BTCUSDT","orderTypes":["LIMIT","MARKET","STOP","STOP_MARKET","TAKE_PROFIT",` +
	`"TAKE_PROFIT_MARKET","TRAILING_STOP_MARKET"],"timeInForce":["GTC","IOC","GTX","HIDDEN"],` +
	`"pair":"BTCUSDT","contractType":"PERPETUAL","status":"TRADING","baseAsset":"BTC","quoteAsset":"USDT",` +
	`"marginAsset":"USDT","pricePrecision":1,"quantityPrecision":3,"triggerProtect":"0.0200","liquidationFee":"0.025000"}`

func TestFuturesSymbolOrderTypes(t *testing.T) {
	var s FuturesSymbol
	if err := common.JSONUnmarshal([]byte(liveSymbol), &s); err != nil {
		t.Fatal(err)
	}
	if len(s.OrderTypes) != 7 || s.OrderTypes[0] != "LIMIT" || s.OrderTypes[6] != "TRAILING_STOP_MARKET" {
		t.Errorf("OrderTypes = %v", s.OrderTypes)
	}
	if len(s.TimeInForce) != 4 {
		t.Errorf("TimeInForce = %v", s.TimeInForce)
	}
}

// liveBookTicker is a live /fapi/v3/ticker/bookTicker?symbol=BTCUSDT response.
const liveBookTicker = `{"symbol":"BTCUSDT","bidPrice":"83303.9","bidQty":"4.797","askPrice":"83304.0",` +
	`"askQty":"0.001","time":1790610988500,"lastUpdateId":571363136052}`

func TestBookTickerLastUpdateID(t *testing.T) {
	var b BookTicker
	if err := common.JSONUnmarshal([]byte(liveBookTicker), &b); err != nil {
		t.Fatal(err)
	}
	if b.LastUpdateId != 571363136052 || !b.Time.Equal(time.UnixMilli(1790610988500)) || b.AskPrice.String() != "83304" {
		t.Errorf("decoded %+v", b)
	}
}

// TestGetRemainingOpenableNotionalValue checks the request and decodes a live
// /fapi/v3/remainingOpenableNotionalValue?symbol=BTCUSDT&leverage=10 response.
func TestGetRemainingOpenableNotionalValue(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"remainingOpenableNotionalValue":"531083199.18960000000"}`))
	}))
	defer srv.Close()

	got, err := NewFuturesClient(client.WithBaseURL(srv.URL)).
		NewGetRemainingOpenableNotionalValueService("BTCUSDT", 10).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotPath != "/fapi/v3/remainingOpenableNotionalValue" {
		t.Errorf("path = %q", gotPath)
	}
	assertCovers(t, gotQuery, map[string]string{"symbol": "BTCUSDT", "leverage": "10"})
	if gotQuery.Has("signature") {
		t.Error("public endpoint was signed")
	}
	if got.RemainingOpenableNotionalValue.String() != "531083199.1896" {
		t.Errorf("RemainingOpenableNotionalValue = %s", got.RemainingOpenableNotionalValue)
	}
}

// liveAssetLogos is the head of a live /fapi/v3/common/asset/all-asset-logo
// response.
const liveAssetLogos = `[{"assetCode":"BTC","logoUrl":"https://static.apollox.com/image/admin_mgs_image_upload/` +
	`20251010/1d7eca3f-2cc2-4d20-b39e-e93e104efc38.png"},{"assetCode":"ETH","logoUrl":"https://static.apollox.com/` +
	`image/admin_mgs_image_upload/20251010/29537d4d-7a9e-4074-baed-60ba2b9a40ce.png"}]`

func TestGetAllAssetLogos(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(liveAssetLogos))
	}))
	defer srv.Close()

	got, err := NewFuturesClient(client.WithBaseURL(srv.URL)).NewGetAllAssetLogosService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotPath != "/fapi/v3/common/asset/all-asset-logo" {
		t.Errorf("path = %q", gotPath)
	}
	if len(got) != 2 || got[1].AssetCode != "ETH" ||
		got[0].LogoURL != "https://static.apollox.com/image/admin_mgs_image_upload/20251010/1d7eca3f-2cc2-4d20-b39e-e93e104efc38.png" {
		t.Errorf("decoded %+v", got)
	}
}

// liveExchangeInfoHead is a live /fapi/v3/exchangeInfo response with the
// symbols and assets dropped.
const liveExchangeInfoHead = `{"timezone":"UTC","serverTime":1791471905124,"futuresType":"U_MARGINED",` +
	`"rateLimits":[{"rateLimitType":"ORDERS","interval":"SECOND","intervalNum":10,"limit":300}],` +
	`"exchangeFilters":[],"assets":[],"symbols":[]}`

func TestExchangeInfoFuturesType(t *testing.T) {
	var e ExchangeInfoResponse
	if err := common.JSONUnmarshal([]byte(liveExchangeInfoHead), &e); err != nil {
		t.Fatal(err)
	}
	if e.FuturesType != "U_MARGINED" || !e.ServerTime.Equal(time.UnixMilli(1791471905124)) ||
		len(e.RateLimits) != 1 || e.RateLimits[0].Interval != "SECOND" {
		t.Errorf("decoded %+v", e)
	}
}
