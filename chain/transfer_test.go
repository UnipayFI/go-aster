package chain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/client"
)

// batchTransferBody is the response sample from the Aster-Chain batch-transfer doc.
const batchTransferBody = `{"clientTranId":"batch-20261026-000000001","status":"I","totalCount":2,` +
	`"items":[{"toAddress":"0xAbCd1234...","asset":"USDT","transferId":"100001"},` +
	`{"toAddress":"0xEf567890...","asset":"USDT","transferId":"100002"}]}`

// batchTransferHistoryBody is the response sample from the Aster-Chain
// batch-transfer/history doc, exercising both succeeded and failed items.
const batchTransferHistoryBody = `{"total":1,"rows":[{"clientTranId":"batch-20261026-000000001",` +
	`"status":"M","createTime":1761000000000,"updateTime":1761000005000,"totalCount":2,` +
	`"successCount":1,"failCount":1,"partialCount":0,"processingCount":0,` +
	`"failedItems":[{"transferId":"100002","toAddress":"0xEf567890...","asset":"USDT","amount":"2",` +
	`"transferChainId":null,"successAmount":null,"status":"F","errorCode":"INSUFFICIENT_BALANCE",` +
	`"errorMsg":"insufficient balance","createTime":1761000000000,"updateTime":1761000005000}],` +
	`"details":[{"transferId":"100001","toAddress":"0xAbCd1234...","asset":"USDT","amount":"1.5",` +
	`"transferChainId":56,"successAmount":"1.5","status":"S","errorCode":null,"errorMsg":null,` +
	`"createTime":1761000000000,"updateTime":1761000003000},` +
	`{"transferId":"100002","toAddress":"0xEf567890...","asset":"USDT","amount":"2",` +
	`"transferChainId":null,"successAmount":null,"status":"F","errorCode":"INSUFFICIENT_BALANCE",` +
	`"errorMsg":"insufficient balance","createTime":1761000000000,"updateTime":1761000005000}]}]}`

// TestBatchTransferRequest locks the path, HTTP method, caller-supplied
// signature parameters and items JSON of POST /aster-chain/v3/batch-transfer,
// and decodes the documented response.
func TestBatchTransferRequest(t *testing.T) {
	var gotMethod, gotPath, gotRawQuery string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotRawQuery, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery, r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(batchTransferBody))
	}))
	defer srv.Close()

	chainID := int64(56)
	items := []BatchTransferItem{
		{ToAddress: "0xAbCd1234...", Asset: "USDT", Amount: "1.5"},
		{ToAddress: "0xEf567890...", Asset: "USDT", Amount: "2", TransferChainID: &chainID},
	}
	wantItems := `[{"toAddress":"0xAbCd1234...","asset":"USDT","amount":"1.5"},` +
		`{"toAddress":"0xEf567890...","asset":"USDT","amount":"2","transferChainId":56}]`

	c := NewChainClient(client.WithBaseURL(srv.URL))
	got, err := c.NewBatchTransferService("batch-20261026-000000001", items, testUserAddress, 1761000000000000, "0xdeadbeef").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if gotMethod != http.MethodPost || gotPath != "/aster-chain/v3/batch-transfer" {
		t.Errorf("request = %s %s, want POST /aster-chain/v3/batch-transfer", gotMethod, gotPath)
	}
	assertCovers(t, gotQuery, map[string]string{
		"clientTranId": "batch-20261026-000000001",
		"items":        wantItems,
		"nonce":        "1761000000000000",
		"user":         testUserAddress,
		"signature":    "0xdeadbeef",
	})
	if gotQuery.Has("signer") {
		t.Errorf("query carries signer = %q, want caller-signed user payload only", gotQuery.Get("signer"))
	}
	wantUnsigned := "clientTranId=batch-20261026-000000001&items=" + url.QueryEscape(wantItems) +
		"&nonce=1761000000000000&user=" + testUserAddress
	if unsigned := strings.Replace(gotRawQuery, "&signature=0xdeadbeef", "", 1); unsigned != wantUnsigned {
		t.Errorf("unsigned query = %q, want %q", unsigned, wantUnsigned)
	}

	if got.ClientTranID != "batch-20261026-000000001" || got.Status != BatchTransferPending || got.TotalCount != 2 || len(got.Items) != 2 {
		t.Fatalf("decoded result = %+v", got)
	}
	if got.Items[0].TransferID != "100001" || got.Items[1].ToAddress != "0xEf567890..." {
		t.Errorf("decoded items = %+v", got.Items)
	}
}

// TestGetBatchTransferHistoryRequest locks the path and signed query of GET
// /aster-chain/v3/batch-transfer/history and checks that both succeeded and
// failed item details decode, including nullable fields and unix-milli times.
func TestGetBatchTransferHistoryRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, batchTransferHistoryBody, &gotPath, &gotQuery)

	start, end := time.UnixMilli(1760900000000), time.UnixMilli(1761100000000)
	got, err := c.NewGetBatchTransferHistoryService().
		SetClientTranId("batch-20261026-000000001").
		SetStatus(BatchTransferPartial).
		SetStartTime(start).
		SetEndTime(end).
		SetPage(1).
		SetLimit(20).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if want := "/aster-chain/v3/batch-transfer/history"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	assertSigned(t, gotQuery, map[string]string{
		"clientTranId": "batch-20261026-000000001",
		"status":       "M",
		"startTime":    "1760900000000",
		"endTime":      "1761100000000",
		"page":         "1",
		"limit":        "20",
	})

	if got.Total != 1 || len(got.Rows) != 1 {
		t.Fatalf("decoded history = %+v", got)
	}
	row := got.Rows[0]
	if row.ClientTranID != "batch-20261026-000000001" || row.Status != BatchTransferPartial ||
		!row.CreateTime.Equal(time.UnixMilli(1761000000000)) || !row.UpdateTime.Equal(time.UnixMilli(1761000005000)) ||
		row.TotalCount != 2 || row.SuccessCount != 1 || row.FailCount != 1 || len(row.FailedItems) != 1 || len(row.Details) != 2 {
		t.Fatalf("decoded row = %+v", row)
	}
	okItem := row.Details[0]
	if okItem.Status != BatchTransferSuccess || okItem.TransferChainID == nil || *okItem.TransferChainID != 56 ||
		okItem.SuccessAmount == nil || okItem.SuccessAmount.String() != "1.5" || okItem.ErrorCode != nil || okItem.ErrorMsg != nil ||
		!okItem.UpdateTime.Equal(time.UnixMilli(1761000003000)) {
		t.Errorf("succeeded item = %+v", okItem)
	}
	failItem := row.FailedItems[0]
	if failItem.Status != BatchTransferFailed || failItem.TransferChainID != nil || failItem.SuccessAmount != nil ||
		failItem.ErrorCode == nil || *failItem.ErrorCode != "INSUFFICIENT_BALANCE" ||
		failItem.ErrorMsg == nil || *failItem.ErrorMsg != "insufficient balance" {
		t.Errorf("failed item = %+v", failItem)
	}
}
