package futures

import (
	"testing"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/shopspring/decimal"
)

func TestBracketAcceptsStringOrNumberInts(t *testing.T) {
	for _, body := range []string{
		`[{"symbol":"BTCUSDT","brackets":[{"bracket":"2","initialLeverage":"50","notionalCap":"250000","notionalFloor":"50000","maintMarginRatio":"0.01","cum":"250"}]}]`,
		`[{"symbol":"BTCUSDT","brackets":[{"bracket":2,"initialLeverage":50,"notionalCap":250000,"notionalFloor":50000,"maintMarginRatio":0.01,"cum":250}]}]`,
	} {
		var got []LeverageBracket
		if err := common.JSONUnmarshal([]byte(body), &got); err != nil {
			t.Fatalf("unmarshal %s: %v", body, err)
		}
		b := got[0].Brackets[0]
		if got[0].Symbol != "BTCUSDT" || b.Bracket != 2 || b.InitialLeverage != 50 ||
			!b.NotionalCap.Equal(decimal.NewFromInt(250000)) || !b.NotionalFloor.Equal(decimal.NewFromInt(50000)) ||
			!b.MaintMarginRatio.Equal(decimal.RequireFromString("0.01")) || !b.Cum.Equal(decimal.NewFromInt(250)) {
			t.Errorf("decoded %+v from %s", got[0], body)
		}
	}
}

func TestBracketRejectsNonNumericInt(t *testing.T) {
	var got []LeverageBracket
	if err := common.JSONUnmarshal([]byte(`[{"symbol":"X","brackets":[{"bracket":"abc"}]}]`), &got); err == nil {
		t.Fatal("want error for non-numeric bracket")
	}
}
