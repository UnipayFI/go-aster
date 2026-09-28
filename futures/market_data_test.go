package futures

import (
	"testing"
	"time"

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
