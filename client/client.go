package client

import (
	"crypto/ecdsa"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
)

// SignFn signs an EIP-712 typed-data payload using the signer's private key.
// The msg argument is the URL-encoded query string of the request parameters
// (placed inside Message.msg). The function returns a hex-encoded signature
// (with the leading "0x" omitted, matching eth_account.sign_message output).
type SignFn = func(privateKeyHex string, msg string, chainID int64) (signature string, err error)

type Client struct {
	client *resty.Client

	userAddress      string
	signerAddress    string
	signerPrivateKey *ecdsa.PrivateKey
	privateKeyHex    string
	chainID          int64
	recvWindow       int64
	logger           log.Logger
	signFn           SignFn
	timeOffsetMs     int64

	authErr error

	nonceMu   sync.Mutex
	lastNonce int64

	UsedWeight UsedWeight
	OrderCount OrderCount
}

type UsedWeight struct {
	Used   int64
	Used1M int64
}

type OrderCount struct {
	Count10s int64
	Count1d  int64
}

func NewClient(options ...Options) *Client {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}

	c := &Client{
		client:        opt.client,
		userAddress:   opt.userAddress,
		chainID:       opt.chainID,
		recvWindow:    opt.recvWindow,
		logger:        opt.logger,
		signFn:        opt.signFn,
		timeOffsetMs:  opt.timeOffsetMs,
		privateKeyHex: opt.signerPrivateKeyHex,
	}

	if opt.signerPrivateKeyHex != "" {
		priv, addr, err := parsePrivateKey(opt.signerPrivateKeyHex)
		if err != nil {
			c.authErr = err
		} else {
			c.signerPrivateKey = priv
			c.signerAddress = addr
		}
	}

	return c
}

func parsePrivateKey(hexKey string) (*ecdsa.PrivateKey, string, error) {
	cleaned := strings.TrimPrefix(strings.TrimSpace(hexKey), "0x")
	if cleaned == "" {
		return nil, "", errors.New("empty private key")
	}
	priv, err := crypto.HexToECDSA(cleaned)
	if err != nil {
		return nil, "", err
	}
	addr := crypto.PubkeyToAddress(priv.PublicKey)
	return priv, common.HexToAddress(addr.Hex()).Hex(), nil
}

func (c *Client) GetHttpClient() *resty.Client {
	return c.client
}

func (c *Client) GetUserAddress() string {
	return c.userAddress
}

func (c *Client) GetSignerAddress() string {
	return c.signerAddress
}

func (c *Client) GetSignerPrivateKey() *ecdsa.PrivateKey {
	return c.signerPrivateKey
}

func (c *Client) GetPrivateKeyHex() string {
	return c.privateKeyHex
}

func (c *Client) GetChainID() int64 {
	return c.chainID
}

func (c *Client) GetRecvWindow() int64 {
	return c.recvWindow
}

func (c *Client) GetLogger() log.Logger {
	return c.logger
}

func (c *Client) GetSignFn() SignFn {
	return c.signFn
}

func (c *Client) GetTimeOffsetMs() int64 {
	return c.timeOffsetMs
}

func (c *Client) SetTimeOffset(offsetMs int64) {
	c.timeOffsetMs = offsetMs
}

func (c *Client) GetAuthError() error {
	return c.authErr
}

// NextNonce returns a microsecond-precision nonce that is monotonically
// increasing across concurrent calls. The local clock is adjusted by
// timeOffsetMs to compensate for client-server time drift.
func (c *Client) NextNonce() int64 {
	c.nonceMu.Lock()
	defer c.nonceMu.Unlock()
	now := time.Now().UnixMicro() - c.timeOffsetMs*1000
	if now <= c.lastNonce {
		now = c.lastNonce + 1
	}
	c.lastNonce = now
	return now
}

func (c *Client) GetUsedWeight() UsedWeight {
	return c.UsedWeight
}

func (c *Client) GetOrderCount() OrderCount {
	return c.OrderCount
}

func (c *Client) SetUsedWeight(used, used1m *int64) {
	if used != nil {
		c.UsedWeight.Used = *used
	}
	if used1m != nil {
		c.UsedWeight.Used1M = *used1m
	}
}

func (c *Client) SetOrderCount(count10s, count1d *int64) {
	if count10s != nil {
		c.OrderCount.Count10s = *count10s
	}
	if count1d != nil {
		c.OrderCount.Count1d = *count1d
	}
}
