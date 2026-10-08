package futures

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/shopspring/decimal"
)

// Reuses the documented V3 demo credentials (see request/sign_test.go) so the
// tests need no environment configuration.
const (
	testUserAddress  = "0x63DD5aCC6b1aa0f563956C0e534DD30B6dcF7C4e"
	testSignerKeyHex = "0x4fd0a42218f3eae43a6ce26d22544e986139a01e5b34a62db53757ffca81bae1"
)

// newTestClient serves body for every request.
func newTestClient(t *testing.T, body string) *FuturesClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return NewFuturesClient(
		client.WithBaseURL(srv.URL),
		client.WithAuth(testUserAddress, testSignerKeyHex),
	)
}

// liveChaseBody is a real POST /fapi/v3/chase response (captured by
// JKorf/Aster.Net, ids redacted): the strategy result, not the echoed order
// parameters the docs show.
const liveChaseBody = `{"strategyId":123,"clientStrategyId":"xyx","strategyType":"CHASE",` +
	`"strategyStatus":"WORKING","updateTime":1779785417130,"failureCode":0,"failureReason":""}`

// docChaseBody is the response sample from the Place Chase Order doc.
const docChaseBody = `{"strategyId":12345,"clientStrategyId":"my_chase_1","symbol":"BTCUSDT",` +
	`"side":"BUY","positionSide":"BOTH","quantity":"0.1","quantityUnit":"BASE","reduceOnly":false,` +
	`"chaseOffset":"0.5","chaseOffsetType":"ABSOLUTE","maxChaseOffset":"10.0",` +
	`"maxChaseOffsetType":"ABSOLUTE","timeInForce":"GTX","strategyStatus":"NEW",` +
	`"bookTime":1747728000000,"updateTime":1747728000000}`

func TestPlaceChaseOrderResponse(t *testing.T) {
	c := newTestClient(t, liveChaseBody)
	got, err := c.NewPlaceChaseOrderService("BTCUSDT", SideBuy, QuantityUnitBase, decimal.NewFromInt(1)).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.StrategyId != 123 || got.StrategyType != StrategyChase || got.StrategyStatus != "WORKING" ||
		got.FailureCode != 0 || got.FailureReason != "" || got.BookTime != nil ||
		!got.UpdateTime.Equal(time.UnixMilli(1779785417130)) {
		t.Errorf("live response decoded as %+v", got)
	}

	c = newTestClient(t, docChaseBody)
	got, err = c.NewPlaceChaseOrderService("BTCUSDT", SideBuy, QuantityUnitBase, decimal.NewFromInt(1)).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.Symbol != "BTCUSDT" || got.ChaseOffset.String() != "0.5" || got.BookTime == nil ||
		!got.BookTime.Equal(time.UnixMilli(1747728000000)) || got.StrategyType != "" {
		t.Errorf("doc response decoded as %+v", got)
	}
}

// docNewOrderBody is the response sample from the New Order doc.
const docNewOrderBody = `{"clientOrderId":"testOrder","cumQty":"0","cumQuote":"0","executedQty":"0",` +
	`"orderId":22542179,"avgPrice":"0.00000","origQty":"10","price":"0","reduceOnly":false,"side":"BUY",` +
	`"positionSide":"SHORT","status":"NEW","stopPrice":"9300","closePosition":false,"symbol":"BTCUSDT",` +
	`"timeInForce":"GTC","type":"TRAILING_STOP_MARKET","origType":"TRAILING_STOP_MARKET","activatePrice":"9020",` +
	`"priceRate":"0.3","updateTime":1566818724722,"workingType":"CONTRACT_PRICE","priceProtect":false,` +
	`"pegPriceType":"","stpMode":""}`

func TestPlaceOrderResponse(t *testing.T) {
	got, err := newTestClient(t, docNewOrderBody).
		NewPlaceOrderService("BTCUSDT", SideBuy, OrderTypeTrailingStopMarket).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.OrderId != 22542179 || got.PegPriceType != "" || got.STPMode != "" || got.Time != nil ||
		!got.UpdateTime.Equal(time.UnixMilli(1566818724722)) {
		t.Errorf("decoded %+v", got)
	}

	body := strings.Replace(docNewOrderBody, `"pegPriceType":"","stpMode":""`,
		`"pegPriceType":"QUEUE_1","stpMode":"EXPIRE_MAKER"`, 1)
	got, err = newTestClient(t, body).
		NewPlaceOrderService("BTCUSDT", SideBuy, OrderTypeLimit).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.PegPriceType != PegQueue1 || got.STPMode != STPExpireMaker {
		t.Errorf("PegPriceType = %q, STPMode = %q", got.PegPriceType, got.STPMode)
	}
}
