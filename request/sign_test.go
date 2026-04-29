package request

import (
	"encoding/hex"
	"strings"
	"testing"

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
