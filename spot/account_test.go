package spot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
)

// liveWithdrawFeeBody is a live /api/v3/aster/withdraw/estimateFee response
// (chainId=56, USDT).
const liveWithdrawFeeBody = `{"gasLimit":200000,"tokenPrice":0.99961,"gasCost":0.11,"gasUsdValue":0.1}`

func TestGetWithdrawFee(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(liveWithdrawFeeBody))
	}))
	defer srv.Close()
	fee, err := NewSpotClient(client.WithBaseURL(srv.URL)).NewGetWithdrawFeeService("56", "USDT").Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if fee.GasLimit != 200000 || fee.TokenPrice.String() != "0.99961" || fee.GasCost.String() != "0.11" || fee.GasUsdValue.String() != "0.1" {
		t.Errorf("decoded %+v", fee)
	}
}
