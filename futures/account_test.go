package futures

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
)

// incomeBody is the response sample from the Get Income History doc, which
// quotes tranId.
const incomeBody = `[{"symbol":"","incomeType":"TRANSFER","income":"-0.37500000","asset":"USDT",` +
	`"info":"TRANSFER","time":1570608000000,"tranId":"9689322392","tradeId":""},` +
	`{"symbol":"BTCUSDT","incomeType":"COMMISSION","income":"-0.01000000","asset":"USDT",` +
	`"info":"COMMISSION","time":1570636800000,"tranId":"9689322393","tradeId":"2059192"}]`

// incomeBodyBare is the same page with tranId as a JSON number.
const incomeBodyBare = `[{"symbol":"","incomeType":"TRANSFER","income":"-0.37500000","asset":"USDT",` +
	`"info":"TRANSFER","time":1570608000000,"tranId":9689322392,"tradeId":""},` +
	`{"symbol":"BTCUSDT","incomeType":"COMMISSION","income":"-0.01000000","asset":"USDT",` +
	`"info":"COMMISSION","time":1570636800000,"tranId":9689322393,"tradeId":"2059192"}]`

// TestIncomeHistoryTranID checks that tranId decodes whether Aster sends it
// quoted (as documented) or as a number, and that the rest of the record,
// including the unixmilli time, still decodes normally.
func TestIncomeHistoryTranID(t *testing.T) {
	for _, body := range []string{incomeBody, incomeBodyBare} {
		got, err := newTestClient(t, body).NewGetIncomeHistoryService().Do(context.Background())
		if err != nil {
			t.Fatalf("Do(%s): %v", body, err)
		}
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
		if got[0].TranID != 9689322392 || got[1].TranID != 9689322393 {
			t.Errorf("TranID = %d, %d, want 9689322392, 9689322393", got[0].TranID, got[1].TranID)
		}
		if !got[1].Time.Equal(time.UnixMilli(1570636800000)) || got[1].Symbol != "BTCUSDT" ||
			got[1].IncomeType != "COMMISSION" || got[1].Income.String() != "-0.01" || got[1].TradeID != "2059192" {
			t.Errorf("record = %+v", got[1])
		}
	}

	for _, body := range []string{`[{"tranId":"abc"}]`, `[{"tranId":1.5}]`, `[{"tranId":true}]`, `{"code":1}`} {
		if _, err := newTestClient(t, body).NewGetIncomeHistoryService().Do(context.Background()); err == nil {
			t.Errorf("%s: want error", body)
		}
	}
}

// TestIncomeHistoryErrors checks that a bad field is reported where it is,
// whether the page quotes tranId or not.
func TestIncomeHistoryErrors(t *testing.T) {
	for body, want := range map[string]string{
		`[{"tranId":"5"},{"tranId":"6","income":"zz"}]`: `/1/income`,
		`[{"tranId":"5","tradeId":7}]`:                  `/0/tradeId`,
		`[{"tranId":5,"time":"bad"}]`:                   `/0/time`,
		`[{"tranId":"5","time":"bad"}]`:                 `/0/time`,
		`[{"tranId":"abc"}]`:                            `/0/tranId`,
		`[{"tranId":5},{"tranId":"x"}]`:                 `/1/tranId`,
		`[{"tranId":"5"},{"tranId":6}]`:                 `/1/tranId`,
	} {
		_, err := newTestClient(t, body).NewGetIncomeHistoryService().Do(context.Background())
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err = %v, want it to point at %s", body, err, want)
		}
	}
}

// TestIncomeRecordQuotedTags keeps incomeRecordQuoted's json tags in step with
// IncomeRecord's; the struct conversion only guards names and types.
func TestIncomeRecordQuotedTags(t *testing.T) {
	a, b := reflect.TypeFor[IncomeRecord](), reflect.TypeFor[incomeRecordQuoted]()
	for i := range a.NumField() {
		want := a.Field(i).Tag.Get("json")
		if a.Field(i).Name == "TranID" {
			want += ",string"
		}
		if got := b.Field(i).Tag.Get("json"); got != want {
			t.Errorf("incomeRecordQuoted.%s json tag = %q, want %q", b.Field(i).Name, got, want)
		}
	}
}

// TestGetAccountPaths checks that the join-margin and plain account services
// call their own paths, signed, and share the response shape.
func TestGetAccountPaths(t *testing.T) {
	for want, do := range map[string]func(*FuturesClient) (*AccountInfo, error){
		"/fapi/v3/accountWithJoinMargin": func(c *FuturesClient) (*AccountInfo, error) {
			return c.NewGetAccountService().Do(context.Background())
		},
		"/fapi/v3/account": func(c *FuturesClient) (*AccountInfo, error) {
			return c.NewGetAccountWithoutJoinMarginService().Do(context.Background())
		},
	} {
		var gotPath string
		var gotQuery url.Values
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath, gotQuery = r.URL.Path, r.URL.Query()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"feeTier":1,"canTrade":true,"totalWalletBalance":"103.12345678"}`))
		}))
		got, err := do(NewFuturesClient(client.WithBaseURL(srv.URL), client.WithAuth(testUserAddress, testSignerKeyHex)))
		srv.Close()
		if err != nil {
			t.Fatalf("%s: %v", want, err)
		}
		if gotPath != want || gotQuery.Get("signature") == "" {
			t.Errorf("path = %q, signature = %q, want %q signed", gotPath, gotQuery.Get("signature"), want)
		}
		if got.FeeTier != 1 || !got.CanTrade || got.TotalWalletBalance.String() != "103.12345678" {
			t.Errorf("%s: decoded %+v", want, got)
		}
	}
}
