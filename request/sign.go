package request

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	asterCommon "github.com/UnipayFI/go-aster/v3/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EIP712Digest computes the 32-byte EIP-712 digest used by Aster V3 signing:
//
//	domainSeparator = keccak256( typeHash(EIP712Domain) || keccak256(name) ||
//	                              keccak256(version) || pad32(chainID) ||
//	                              pad32(verifyingContract) )
//	structHash      = keccak256( typeHash(Message) || keccak256(msg) )
//	digest          = keccak256( 0x1901 || domainSeparator || structHash )
//
// The domain values are fixed by the protocol (name="AsterSignTransaction",
// version="1", verifyingContract=0x000...0); only chainID is configurable.
func EIP712Digest(msg string, chainID int64) []byte {
	domainTypeHash := crypto.Keccak256([]byte(
		"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)",
	))
	nameHash := crypto.Keccak256([]byte(asterCommon.EIP712_DOMAIN_NAME))
	versionHash := crypto.Keccak256([]byte(asterCommon.EIP712_DOMAIN_VERSION))
	chainIDPadded := leftPad32(big.NewInt(chainID).Bytes())
	verifyingContract := ethCommon.HexToAddress(asterCommon.EIP712_VERIFYING_CONTRACT)
	verifyingContractPadded := leftPad32(verifyingContract.Bytes())

	domainSeparator := crypto.Keccak256(
		domainTypeHash,
		nameHash,
		versionHash,
		chainIDPadded,
		verifyingContractPadded,
	)

	messageTypeHash := crypto.Keccak256([]byte("Message(string msg)"))
	messageValueHash := crypto.Keccak256([]byte(msg))
	structHash := crypto.Keccak256(messageTypeHash, messageValueHash)

	return crypto.Keccak256([]byte{0x19, 0x01}, domainSeparator, structHash)
}

// SignEIP712V3 signs a Message-typed EIP-712 payload with the given private
// key and returns a 65-byte hex signature with v adjusted to 27/28 (matching
// eth_account.sign_message in the official Python demo).
func SignEIP712V3(privateKeyHex, msg string, chainID int64) (string, error) {
	priv, err := parsePrivateKeyHex(privateKeyHex)
	if err != nil {
		return "", err
	}
	return signWithKey(priv, msg, chainID)
}

func signWithKey(priv *ecdsa.PrivateKey, msg string, chainID int64) (string, error) {
	digest := EIP712Digest(msg, chainID)
	sig, err := crypto.Sign(digest, priv)
	if err != nil {
		return "", err
	}
	if len(sig) != 65 {
		return "", fmt.Errorf("unexpected signature length: %d", len(sig))
	}
	sig[64] += 27
	return hex.EncodeToString(sig), nil
}

func parsePrivateKeyHex(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	cleaned := strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	if cleaned == "" {
		return nil, errors.New("empty private key")
	}
	return crypto.HexToECDSA(cleaned)
}

func leftPad32(b []byte) []byte {
	if len(b) >= 32 {
		return b[len(b)-32:]
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
