package chain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
)

// Reuses the documented V3 demo credentials (see request/sign_test.go) so the
// test needs no environment configuration.
const (
	testUserAddress   = "0x63DD5aCC6b1aa0f563956C0e534DD30B6dcF7C4e"
	testSignerKeyHex  = "0x4fd0a42218f3eae43a6ce26d22544e986139a01e5b34a62db53757ffca81bae1"
	testSignerAddress = "0x21cF8Ae13Bb72632562c6Fff438652Ba1a151bb0"
)

// withdrawInfoBody is the response sample from the Aster-Chain endpoints doc.
const withdrawInfoBody = `{"userDailyLimit":"10000","userRemainingDailyLimit":"9500",` +
	`"totalDailyLimit":"100000","totalRemainingDailyLimit":"95000","balances":{"USDT":{` +
	`"currency":"USDT","spotTotalWithdrawAmount":"0","perpTotalWithdrawAmount":"500",` +
	`"dailyLimit":"5000","chainBalances":{"1":{"chainId":1,"spotMaxWithdrawAmount":"1000",` +
	`"perpMaxWithdrawAmount":"4500","chainLimit":"5000","withdrawFee":"0.5"}}}}}`

// historyBody is the response sample from the Aster-Chain endpoints doc.
const historyBody = `[{"id":"12345","type":"WITHDRAW","asset":"USDT","amount":"100",` +
	`"state":"COMPLETED","txHash":"0xabc123...","time":1699900800000,"chainId":1,` +
	`"accountType":"perp"}]`

// estimateFeeBody is the response sample from the Aster-Chain endpoints doc.
// Note gasPrice and gasLimit arrive as JSON numbers, the rest as strings.
const estimateFeeBody = `{"gasPrice":1000000000,"gasLimit":21000,"nativePrice":"1800.00",` +
	`"tokenPrice":"1.00","gasCost":"0.000021","gasUsdValue":"0.038"}`

// depositAddressBody is the response sample from the Aster-Chain endpoints doc.
const depositAddressBody = `{"network":"SUI",` +
	`"address":"0x9a40f0119b670fb6b155744b51981f91c4c4c8a20c333441a63853fe7d055c90"}`

// newTestClient serves body for every request and records the request line.
func newTestClient(t *testing.T, body string, gotPath *string, gotQuery *url.Values) *ChainClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*gotPath = r.URL.Path
		*gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return NewChainClient(
		client.WithBaseURL(srv.URL),
		client.WithAuth(testUserAddress, testSignerKeyHex),
	)
}

// TestGetPerpWithdrawInfoRequest locks the request shape of GET
// /aster-chain/v3/perp/user-withdraw-info -- the path plus the four signed
// query parameters -- and checks that the documented response decodes,
// including the asset- and chain-keyed nested maps.
func TestGetPerpWithdrawInfoRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, withdrawInfoBody, &gotPath, &gotQuery)

	info, err := c.NewGetPerpWithdrawInfoService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if want := "/aster-chain/v3/perp/user-withdraw-info"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	assertSigned(t, gotQuery, nil)

	if want := "9500"; info.UserRemainingDailyLimit.String() != want {
		t.Errorf("UserRemainingDailyLimit = %s, want %s", info.UserRemainingDailyLimit, want)
	}
	usdt, ok := info.Balances["USDT"]
	if !ok {
		t.Fatalf("Balances is missing the USDT entry: %v", info.Balances)
	}
	if want := "500"; usdt.PerpTotalWithdrawAmount.String() != want {
		t.Errorf("PerpTotalWithdrawAmount = %s, want %s", usdt.PerpTotalWithdrawAmount, want)
	}
	eth, ok := usdt.ChainBalances["1"]
	if !ok {
		t.Fatalf("ChainBalances is missing the \"1\" entry: %v", usdt.ChainBalances)
	}
	if eth.ChainID != 1 {
		t.Errorf("ChainID = %d, want 1", eth.ChainID)
	}
	if want := "0.5"; eth.WithdrawFee.String() != want {
		t.Errorf("WithdrawFee = %s, want %s", eth.WithdrawFee, want)
	}
}

// TestGetDepositWithdrawHistoryRequest locks the request shape of GET
// /aster-chain/v3/perp/deposit-withdraw-history, including the optional
// chainId filter, and checks that time decodes from unix millis.
func TestGetDepositWithdrawHistoryRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, historyBody, &gotPath, &gotQuery)

	records, err := c.NewGetDepositWithdrawHistoryService().SetChainId(1).Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if want := "/aster-chain/v3/perp/deposit-withdraw-history"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	assertSigned(t, gotQuery, map[string]string{"chainId": "1"})

	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}
	got := records[0]
	if got.Type != "WITHDRAW" {
		t.Errorf("Type = %q, want %q", got.Type, "WITHDRAW")
	}
	if got.TxHash != "0xabc123..." {
		t.Errorf("TxHash = %q", got.TxHash)
	}
	if got.Time != 1699900800000 {
		t.Errorf("Time = %d, want 1699900800000", got.Time)
	}
	if got.TimeAt().UnixMilli() != 1699900800000 {
		t.Errorf("TimeAt() = %v, want unix milli 1699900800000", got.TimeAt())
	}
}

// TestGetSpotDepositAddressRequest locks the request shape of GET
// /aster-chain/v3/spot/user-deposit-address, including the optional network
// parameter.
func TestGetSpotDepositAddressRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, depositAddressBody, &gotPath, &gotQuery)

	addr, err := c.NewGetSpotDepositAddressService().SetNetwork("SUI").Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if want := "/aster-chain/v3/spot/user-deposit-address"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	assertSigned(t, gotQuery, map[string]string{"network": "SUI"})

	if addr.Network != "SUI" {
		t.Errorf("Network = %q, want %q", addr.Network, "SUI")
	}
	want := "0x9a40f0119b670fb6b155744b51981f91c4c4c8a20c333441a63853fe7d055c90"
	if addr.Address != want {
		t.Errorf("Address = %q, want %q", addr.Address, want)
	}
}

// TestEstimateWithdrawFeeRequest locks the request shape of GET
// /aster-chain/v3/withdraw/estimateFee. The endpoint is public, so the request
// must carry only the two documented parameters and no signature.
func TestEstimateWithdrawFeeRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, estimateFeeBody, &gotPath, &gotQuery)

	fee, err := c.NewEstimateWithdrawFeeService(1, "USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if want := "/aster-chain/v3/withdraw/estimateFee"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	assertCovers(t, gotQuery, map[string]string{"chainId": "1", "asset": "USDT"})
	for _, key := range []string{"user", "signer", "nonce", "signature"} {
		if gotQuery.Get(key) != "" {
			t.Errorf("query %q = %q, want it absent on a public endpoint", key, gotQuery.Get(key))
		}
	}

	if fee.GasPrice != 1000000000 {
		t.Errorf("GasPrice = %d, want 1000000000", fee.GasPrice)
	}
	if fee.GasLimit != 21000 {
		t.Errorf("GasLimit = %d, want 21000", fee.GasLimit)
	}
	if want := "0.000021"; fee.GasCost.String() != want {
		t.Errorf("GasCost = %s, want %s", fee.GasCost, want)
	}
}

// assertSigned checks that the query carries the V3 signature parameters plus
// every expected endpoint parameter.
func assertSigned(t *testing.T, got url.Values, want map[string]string) {
	t.Helper()
	assertCovers(t, got, map[string]string{
		"user":   testUserAddress,
		"signer": testSignerAddress,
	})
	for _, key := range []string{"nonce", "signature"} {
		if got.Get(key) == "" {
			t.Errorf("query %q is empty, want it signed in", key)
		}
	}
	assertCovers(t, got, want)
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
