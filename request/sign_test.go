package request

import (
	"context"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/ethereum/go-ethereum/crypto"
)

// Test vector taken from
// /tmp/api-docs/V3(Recommended)/EN/aster-finance-spot-api-v3.md (lines
// 121-128) and aster-finance-futures-api-v3.md (lines 244-248):
//
//	user       0x63DD5aCC6b1aa0f563956C0e534DD30B6dcF7C4e
//	signer     0x21cF8Ae13Bb72632562c6Fff438652Ba1a151bb0
//	privateKey 0x4fd0a42218f3eae43a6ce26d22544e986139a01e5b34a62db53757ffca81bae1
//
// Signing uses primaryType="Message" with a single string field "msg" whose
// value is the URL-encoded query string of the request parameters (preserving
// insertion order; nonce and signer appended at the end).
const (
	testPrivateKeyHex = "0x4fd0a42218f3eae43a6ce26d22544e986139a01e5b34a62db53757ffca81bae1"
	testSignerAddress = "0x21cF8Ae13Bb72632562c6Fff438652Ba1a151bb0"
	testChainID       = int64(1666)
	// Mirrors the place_order example in the spot V3 doc:
	//   {"symbol": "ASTERUSDT", "type": "LIMIT", "side": "BUY",
	//    "timeInForce": "GTC", "quantity": "100", "price": "0.4"}
	// after appending nonce and signer (Python urllib.parse.urlencode order).
	testMsg = "symbol=ASTERUSDT&type=LIMIT&side=BUY&timeInForce=GTC&quantity=100&price=0.4&nonce=1748310859508867&signer=0x21cF8Ae13Bb72632562c6Fff438652Ba1a151bb0"
)

// TestEIP712DigestStable locks the digest output so future refactors of the
// hashing path (or upstream go-ethereum changes) cannot silently shift it.
func TestEIP712DigestStable(t *testing.T) {
	digest := EIP712Digest(testMsg, testChainID)
	if len(digest) != 32 {
		t.Fatalf("digest length = %d, want 32", len(digest))
	}
	got := hex.EncodeToString(digest)
	t.Logf("EIP-712 digest: 0x%s", got)
}

// TestSignAndRecover signs the test message and verifies that ecrecover on the
// resulting signature returns the expected signer address. This is the
// strongest self-consistent check we can run without comparing byte-for-byte
// against a Python eth_account output: if the EIP-712 type strings, domain
// values, byte ordering, or v adjustment were wrong, recovery would produce a
// different address.
func TestSignAndRecover(t *testing.T) {
	signatureHex, err := SignEIP712V3(testPrivateKeyHex, testMsg, testChainID)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	if len(sig) != 65 {
		t.Fatalf("signature length = %d, want 65", len(sig))
	}
	if sig[64] != 27 && sig[64] != 28 {
		t.Fatalf("signature v = %d, want 27 or 28 (eth_account format)", sig[64])
	}

	// Undo the v += 27 adjustment for ecrecover (which expects 0/1).
	recoverable := make([]byte, 65)
	copy(recoverable, sig)
	recoverable[64] -= 27

	digest := EIP712Digest(testMsg, testChainID)
	pub, err := crypto.SigToPub(digest, recoverable)
	if err != nil {
		t.Fatalf("recover pubkey: %v", err)
	}
	got := crypto.PubkeyToAddress(*pub).Hex()

	if !strings.EqualFold(got, testSignerAddress) {
		t.Fatalf("recovered signer = %s, want %s", got, testSignerAddress)
	}
}

// TestSignDeterministic confirms that signing the same payload twice yields
// identical output (Aster's signing scheme uses RFC-6979 deterministic ECDSA
// via go-ethereum/crypto, matching eth_account's behavior).
func TestSignDeterministic(t *testing.T) {
	first, err := SignEIP712V3(testPrivateKeyHex, testMsg, testChainID)
	if err != nil {
		t.Fatalf("sign 1: %v", err)
	}
	second, err := SignEIP712V3(testPrivateKeyHex, testMsg, testChainID)
	if err != nil {
		t.Fatalf("sign 2: %v", err)
	}
	if first != second {
		t.Fatalf("signatures differ:\n  1st: %s\n  2nd: %s", first, second)
	}
}

// TestDifferentChainIDProducesDifferentSignature ensures chainID actually
// reaches the digest (a regression here would silently let testnet signatures
// validate on mainnet and vice versa).
func TestDifferentChainIDProducesDifferentSignature(t *testing.T) {
	mainnet, err := SignEIP712V3(testPrivateKeyHex, testMsg, 1666)
	if err != nil {
		t.Fatal(err)
	}
	testnet, err := SignEIP712V3(testPrivateKeyHex, testMsg, 714)
	if err != nil {
		t.Fatal(err)
	}
	if mainnet == testnet {
		t.Fatalf("mainnet and testnet signatures are identical: %s", mainnet)
	}
}

// TestNonceWindowSigning checks client-level and per-request nonceWindow
// handling, including suppression, guarded-cancel isolation, and EIP-712
// coverage of nonceWindow in the signed message.
func TestNonceWindowSigning(t *testing.T) {
	const userAddr = "0x63DD5aCC6b1aa0f563956C0e534DD30B6dcF7C4e"
	ctx := context.Background()

	assertValidSig := func(t *testing.T, rawURL string) url.Values {
		t.Helper()
		rawQuery := strings.SplitN(rawURL, "?", 2)[1]
		idx := strings.LastIndex(rawQuery, "&signature=")
		if idx < 0 {
			t.Fatalf("missing trailing signature in %q", rawQuery)
		}
		unsignedMsg := rawQuery[:idx]
		sigHex := rawQuery[idx+len("&signature="):]
		sig, err := hex.DecodeString(sigHex)
		if err != nil || len(sig) != 65 {
			t.Fatalf("decode signature %q: %v", sigHex, err)
		}
		sig[64] -= 27
		pub, err := crypto.SigToPub(EIP712Digest(unsignedMsg, 1666), sig)
		if err != nil {
			t.Fatalf("ecrecover(%q): %v", unsignedMsg, err)
		}
		if got := crypto.PubkeyToAddress(*pub).Hex(); !strings.EqualFold(got, testSignerAddress) {
			t.Fatalf("recovered signer = %s, want %s", got, testSignerAddress)
		}
		q, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatal(err)
		}
		return q
	}

	// Default client omits nonceWindow.
	plain := client.NewClient(client.ProductFutures, client.WithAuth(userAddr, testPrivateKeyHex))
	q := assertValidSig(t, Post(plain, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT"}).WithSignature().fullURL())
	if q.Has("nonceWindow") {
		t.Errorf("default client sent nonceWindow = %v", q["nonceWindow"])
	}

	// Client-wide WithNonceWindow applies to WithSignature, can be overridden or
	// cleared per request, and does not leak onto guarded cancels.
	withDefault := client.NewClient(client.ProductFutures, client.WithAuth(userAddr, testPrivateKeyHex), client.WithNonceWindow(5000))
	q = assertValidSig(t, Post(withDefault, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT"}).WithSignature().fullURL())
	if got := q["nonceWindow"]; len(got) != 1 || got[0] != "5000" {
		t.Errorf("client default nonceWindow = %v, want [5000]", got)
	}

	q = assertValidSig(t, Post(withDefault, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT", "nonceWindow": "1200"}).WithSignature().fullURL())
	if got := q["nonceWindow"]; len(got) != 1 || got[0] != "1200" {
		t.Errorf("per-param override nonceWindow = %v, want [1200]", got)
	}

	q = assertValidSig(t, Post(withDefault, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT"}).WithNonceWindow(800).WithSignature().fullURL())
	if got := q["nonceWindow"]; len(got) != 1 || got[0] != "800" {
		t.Errorf("WithNonceWindow override = %v, want [800]", got)
	}

	q = assertValidSig(t, Post(withDefault, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT", "nonceWindow": ""}).WithSignature().fullURL())
	if q.Has("nonceWindow") {
		t.Errorf("suppressed param sent nonceWindow = %v", q["nonceWindow"])
	}

	q = assertValidSig(t, Post(withDefault, ctx, "/fapi/v3/order", map[string]string{"symbol": "BTCUSDT"}).WithNonceWindow(0).WithSignature().fullURL())
	if q.Has("nonceWindow") {
		t.Errorf("WithNonceWindow(0) sent nonceWindow = %v", q["nonceWindow"])
	}

	q = assertValidSig(t, Delete(ctx, withDefault, "/fapi/v3/guardedCancelOrder", map[string]string{"symbol": "BTCUSDT", "orderId": "1"}).WithSignatureNonce(1748310859508867).fullURL())
	if q.Has("nonceWindow") {
		t.Errorf("guarded cancel inherited client nonceWindow = %v", q["nonceWindow"])
	}

	q = assertValidSig(t, Delete(ctx, withDefault, "/fapi/v3/guardedCancelOrder", map[string]string{"symbol": "BTCUSDT", "orderId": "1", "nonceWindow": "10000"}).WithSignatureNonce(1748310859508867).fullURL())
	if got := q["nonceWindow"]; len(got) != 1 || got[0] != "10000" {
		t.Errorf("explicit guarded cancel nonceWindow = %v, want [10000]", got)
	}
}
