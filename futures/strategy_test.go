package futures

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
)

// TestStrategyOrderPeg checks that a BBO-pegged leg sends pegPriceType in
// subOrderList while an unpegged leg leaves it out.
func TestStrategyOrderPeg(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"strategyId":123456789,"strategyType":"OTO","strategyStatus":"WORKING",` +
			`"updateTime":1699900000000,"failureCode":-4142,"failureReason":"would trigger immediately"}`))
	}))
	defer srv.Close()
	c := NewFuturesClient(client.WithBaseURL(srv.URL), client.WithAuth(testUserAddress, testSignerKeyHex))

	got, err := c.NewPlaceStrategyOrderService(StrategyOTO, []StrategySubOrder{
		{StrategySubId: "1", SecurityType: "USDT_FUTURES", Symbol: "BTCUSDT", Side: "BUY", Type: "LIMIT",
			Quantity: "0.001", Price: "0", PegPriceType: string(PegQueue5), TimeInForce: "GTC"},
		{StrategySubId: "2", SecurityType: "USDT_FUTURES", Symbol: "BTCUSDT", Side: "SELL", Type: "TAKE_PROFIT_MARKET",
			StopPrice: "60000", ClosePosition: "true", FirstDrivenId: "1", FirstDrivenOn: "FILLED", FirstTrigger: "PLACE_ORDER"},
	}).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got.FailureCode != -4142 {
		t.Errorf("FailureCode = %d, want -4142", got.FailureCode)
	}
	legs := gotQuery.Get("subOrderList")
	if !strings.Contains(legs, `"price":"0","pegPriceType":"QUEUE_5"`) || strings.Count(legs, "pegPriceType") != 1 {
		t.Errorf("subOrderList = %s, want pegPriceType on the first leg only", legs)
	}
}

// TestStrategyOrderDetailPeg checks that a BBO leg's pegPriceType decodes and
// that a non-peg leg, which omits the key, leaves it empty.
func TestStrategyOrderDetailPeg(t *testing.T) {
	body := `{"strategyId":123456789,"strategyType":"OTO","strategyStatus":"WORKING","bookTime":1699900000000,` +
		`"updateTime":1699900001000,"subOrders":[` +
		`{"strategySubId":1,"type":"LIMIT","price":"0","pegPriceType":"QUEUE_1","triggerTime":0},` +
		`{"strategySubId":2,"type":"TAKE_PROFIT","price":"60000","triggerTime":0}]}`
	got, err := newTestClient(t, body).NewGetStrategyOpenOrderService(StrategyOTO).SetStrategyId(123456789).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(got.SubOrders) != 2 || got.SubOrders[0].PegPriceType != PegQueue1 || !got.SubOrders[0].Price.IsZero() ||
		got.SubOrders[1].PegPriceType != "" || got.SubOrders[1].Price.String() != "60000" {
		t.Errorf("SubOrders = %+v", got.SubOrders)
	}
}
