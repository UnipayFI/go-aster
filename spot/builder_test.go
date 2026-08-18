package spot

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

// builderListBody is the response sample from the Spot Aster Code endpoints
// doc. Note maxFeeRate arrives as a JSON number, not a string.
const builderListBody = `[{"userAddress":"0x63DD5aCC6b1aa0f563956C0e534DD30B6dcF7C4e",` +
	`"builderAddress":"0xYourBuilderAddress","maxFeeRate":0.00001,"builderName":"myBuilder"}]`

// TestGetBuildersRequest locks the request shape of GET /api/v3/builder -- the
// path plus the four signed query parameters -- and checks that the documented
// response decodes, including maxFeeRate arriving as a bare JSON number.
func TestGetBuildersRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(builderListBody))
	}))
	defer srv.Close()

	c := NewSpotClient(
		client.WithBaseURL(srv.URL),
		client.WithAuth(testUserAddress, testSignerKeyHex),
	)

	builders, err := c.NewGetBuildersService().Do(context.Background())
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if gotPath != "/api/v3/builder" {
		t.Errorf("path = %q, want %q", gotPath, "/api/v3/builder")
	}
	assertCovers(t, gotQuery, map[string]string{
		"user":   testUserAddress,
		"signer": testSignerAddress,
	})
	for _, key := range []string{"nonce", "signature"} {
		if gotQuery.Get(key) == "" {
			t.Errorf("query %q is empty, want it signed in", key)
		}
	}

	if len(builders) != 1 {
		t.Fatalf("len(builders) = %d, want 1", len(builders))
	}
	got := builders[0]
	if got.BuilderAddress != "0xYourBuilderAddress" {
		t.Errorf("BuilderAddress = %q", got.BuilderAddress)
	}
	if got.BuilderName != "myBuilder" {
		t.Errorf("BuilderName = %q", got.BuilderName)
	}
	if want := "0.00001"; got.MaxFeeRate.String() != want {
		t.Errorf("MaxFeeRate = %s, want %s", got.MaxFeeRate, want)
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
