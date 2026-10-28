package futures

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
)

// builderPage wraps rows in the pagination envelope the builder user queries
// share, with one errors entry.
func builderPage(rows string) string {
	return `{"total":1,"currentPage":1,"totalPages":1,"pageSize":50,"hasMore":false,"rows":[` + rows + `],` +
		`"errors":[{"address":"0xdead","errorMsg":"This user address is not bound to this builder."}]}`
}

// The row samples below are from the Builder query endpoint docs.
const (
	builderAccountRow = `{"address":"0x1234...abcd","feeTier":0,"canTrade":true,"canDeposit":true,` +
		`"canWithdraw":true,"updateTime":1751500000000,"accountType":0,"dualSidePosition":false,` +
		`"jointMargin":false,"feeBurn":false,"feeBurnAssetId":0,` +
		`"symbolConfig":[{"symbol":"BTCUSDT","leverage":20,"notionalLimitCoef":"10"}]}`
	builderOpenOrderRow = `{"address":"0x1234...abcd","symbol":"BTCUSDT","orderId":1917641,"clientOrderId":"abc",` +
		`"price":"0","origQty":"0.40","executedQty":"0","avgPrice":"0.00000","stopPrice":"9300","status":"NEW",` +
		`"side":"BUY","positionSide":"SHORT","type":"TRAILING_STOP_MARKET","timeInForce":"GTC",` +
		`"time":1579276756075,"workingType":"CONTRACT_PRICE"}`
	builderBalanceRow = `{"address":"0x1234...abcd","asset":"USDT","walletBalance":"23.72469206","price":"1",` +
		`"balanceInUsd":"23.72469206"}`
	builderPositionRow = `{"address":"0x1234...abcd","symbol":"BTCUSDT","positionAmt":"20.000",` +
		`"entryPrice":"6563.66500","markPrice":"6679.50671178","unRealizedProfit":"2316.83423560",` +
		`"liquidationPrice":"5930.78","leverage":"10","maxNotionalValue":"20000000","marginType":"isolated",` +
		`"isolatedMargin":"15517.54150468","isAutoAddMargin":"false","positionSide":"LONG",` +
		`"notional":"133593.13423560","isolatedWallet":"13200.70726908","updateTime":1625474304765}`
	builderCommissionRow = `{"address":"0x1234...abcd","symbol":"BTCUSDT","makerCommissionRate":"0.0002",` +
		`"takerCommissionRate":"0.0004"}`
	builderTradeRow = `{"tradeId":1,"insertTime":1751500000000,"symbol":"BTCUSDT","positionSide":"BOTH",` +
		`"price":"78190.5","qty":"0.0002","side":"SELL","fee":"-0.07819010","userAddress":"0x1234...abcd",` +
		`"builderFee":"0.01"}`
	builderFeeRow = `{"tradeId":698759,"insertTime":1751500000000,"orderId":25851813,"symbol":"BTCUSDT",` +
		`"feeAsset":"USDT","marginAsset":"USDT","builderFee":"0.00050000","userAddress":"0x1234...abcd"}`
	builderOrderRow = `{"avgPrice":"0.00000","clientOrderId":"abc","cumQuote":"0","executedQty":"0",` +
		`"orderId":1917641,"origQty":"0.40","origType":"TRAILING_STOP_MARKET","price":"0","reduceOnly":false,` +
		`"side":"BUY","positionSide":"SHORT","status":"NEW","stopPrice":"9300","closePosition":false,` +
		`"symbol":"BTCUSDT","time":1579276756075,"timeInForce":"GTC","type":"TRAILING_STOP_MARKET",` +
		`"activatePrice":"9020","priceRate":"0.3","updateTime":1579276756075,"workingType":"CONTRACT_PRICE",` +
		`"priceProtect":false,"address":"0x1234...abcd"}`
)

// TestBuilderUserQueries locks the path and query of each builder user query
// and checks that its documented response decodes, errors included.
func TestBuilderUserQueries(t *testing.T) {
	start, end := time.UnixMilli(1751400000000), time.UnixMilli(1751500000000)
	for _, tt := range []struct {
		path  string
		row   string
		do    func(*FuturesClient) (any, error)
		want  map[string]string
		check func(t *testing.T, got any)
	}{
		{
			path: "/fapi/v3/builder/userAccounts",
			row:  builderAccountRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserAccountsService().SetUserAddresses("0xa", "0xb").SetSymbol("BTCUSDT").
					SetPage(2).SetLimit(100).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa,0xb", "symbol": "BTCUSDT", "page": "2", "limit": "100"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserAccounts)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.Address != "0x1234...abcd" || !r.CanTrade || !r.UpdateTime.Equal(time.UnixMilli(1751500000000)) ||
					len(r.SymbolConfig) != 1 || r.SymbolConfig[0].Leverage != 20 || r.SymbolConfig[0].NotionalLimitCoef.String() != "10" {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userOpenOrders",
			row:  builderOpenOrderRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserOpenOrdersService().SetUserAddresses("0xa").SetSymbol("BTCUSDT").
					SetPage(1).SetLimit(50).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa", "symbol": "BTCUSDT", "page": "1", "limit": "50"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserOpenOrders)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.OrderID != 1917641 || r.ClientOrderID != "abc" || r.StopPrice.String() != "9300" ||
					r.Type != OrderTypeTrailingStopMarket || !r.Time.Equal(time.UnixMilli(1579276756075)) {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userBalances",
			row:  builderBalanceRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserBalancesService().SetUserAddresses("0xa").SetPage(3).SetLimit(10).
					Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa", "page": "3", "limit": "10"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserBalances)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.Asset != "USDT" || r.BalanceInUSD.String() != "23.72469206" || r.Price.String() != "1" {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userPositionRisk",
			row:  builderPositionRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserPositionRiskService().SetUserAddresses("0xa").SetSymbol("BTCUSDT").
					SetPage(1).SetLimit(50).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa", "symbol": "BTCUSDT", "page": "1", "limit": "50"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserPositionRisk)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.PositionSide != PositionSideLong || r.Notional.String() != "133593.1342356" ||
					r.IsolatedWallet.String() != "13200.70726908" || r.MarginType != "isolated" ||
					!r.UpdateTime.Equal(time.UnixMilli(1625474304765)) {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userCommissionRates",
			row:  builderCommissionRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserCommissionRatesService("BTCUSDT").SetUserAddresses("0xa").
					SetPage(1).SetLimit(50).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa", "symbol": "BTCUSDT", "page": "1", "limit": "50"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserCommissionRates)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.MakerCommissionRate.String() != "0.0002" || r.TakerCommissionRate.String() != "0.0004" {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userTrades",
			row:  builderTradeRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserTradesService().SetUserAddresses("0xa", "0xb").SetStartTime(start).
					SetEndTime(end).SetPage(1).SetLimit(50).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa,0xb", "startTime": "1751400000000",
				"endTime": "1751500000000", "page": "1", "limit": "50"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserTrades)
				checkBuilderErrors(t, p.Errors)
				if r := p.Rows[0]; r.Fee.String() != "-0.0781901" || r.UserAddress != "0x1234...abcd" {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userCommissions",
			row:  builderFeeRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserCommissionsService().SetUserAddresses("0xa", "0xb").SetStartTime(start).
					SetEndTime(end).SetPage(2).SetLimit(1000).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa,0xb", "startTime": "1751400000000",
				"endTime": "1751500000000", "page": "2", "limit": "1000"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserCommissions)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.TradeID != 698759 || r.OrderID != 25851813 || r.MarginAsset != "USDT" ||
					r.BuilderFee.String() != "0.0005" || r.UserAddress != "0x1234...abcd" ||
					!r.InsertTime.Equal(time.UnixMilli(1751500000000)) {
					t.Errorf("row = %+v", r)
				}
			},
		},
		{
			path: "/fapi/v3/builder/userAllOrders",
			row:  builderOrderRow,
			do: func(c *FuturesClient) (any, error) {
				return c.NewGetBuilderUserAllOrdersService().SetUserAddresses("0xa").SetSymbol("BTCUSDT").
					SetStartTime(start).SetEndTime(end).SetPage(1).SetLimit(50).Do(context.Background())
			},
			want: map[string]string{"userAddresses": "0xa", "symbol": "BTCUSDT", "startTime": "1751400000000",
				"endTime": "1751500000000", "page": "1", "limit": "50"},
			check: func(t *testing.T, got any) {
				p := got.(*BuilderUserOrders)
				checkBuilderErrors(t, p.Errors)
				r := p.Rows[0]
				if r.Address != "0x1234...abcd" || r.OrderID != 1917641 || r.ActivatePrice.String() != "9020" ||
					r.OrigType != OrderTypeTrailingStopMarket || !r.UpdateTime.Equal(time.UnixMilli(1579276756075)) {
					t.Errorf("row = %+v", r)
				}
			},
		},
	} {
		t.Run(tt.path, func(t *testing.T) {
			var gotPath string
			var gotQuery url.Values
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotQuery = r.URL.Path, r.URL.Query()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(builderPage(tt.row)))
			}))
			defer srv.Close()
			c := NewFuturesClient(client.WithBaseURL(srv.URL), client.WithAuth(testUserAddress, testSignerKeyHex))

			got, err := tt.do(c)
			if err != nil {
				t.Fatalf("Do: %v", err)
			}
			if gotPath != tt.path {
				t.Errorf("path = %q, want %q", gotPath, tt.path)
			}
			assertCovers(t, gotQuery, tt.want)
			for _, key := range []string{"nonce", "signer", "signature"} {
				if gotQuery.Get(key) == "" {
					t.Errorf("query %q is empty, want it signed in", key)
				}
			}
			tt.check(t, got)
		})
	}
}

func checkBuilderErrors(t *testing.T, errs []BuilderUserError) {
	t.Helper()
	if len(errs) != 1 || errs[0].Address != "0xdead" || errs[0].ErrorMsg != "This user address is not bound to this builder." {
		t.Errorf("Errors = %+v", errs)
	}
}

// assertCovers checks that every expected parameter is present in the query
// with the expected value; extra parameters are ignored.
func assertCovers(t *testing.T, got url.Values, want map[string]string) {
	t.Helper()
	for key, value := range want {
		if got.Get(key) != value {
			t.Errorf("query %q = %q, want %q", key, got.Get(key), value)
		}
	}
}
