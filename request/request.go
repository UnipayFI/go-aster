package request

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/UnipayFI/go-aster/v3/common"
	"github.com/UnipayFI/go-aster/v3/pkg/log"
	"github.com/go-resty/resty/v2"
)

// SignFn signs an EIP-712 typed-data payload using the signer's private key.
// msg is the URL-encoded query string of the request (already including the
// nonce and signer fields). Returns a hex-encoded 65-byte signature with v
// adjusted to 27/28 (matches eth_account.sign_message in the Aster Python demo).
type SignFn = func(privateKeyHex string, msg string, chainID int64) (signature string, err error)

// Client is what every endpoint Service expects from a Spot/Futures client.
// All getters are read-only; SetUsedWeight/SetOrderCount are called from Do
// after parsing rate-limit response headers.
type Client interface {
	GetHttpClient() *resty.Client
	GetUserAddress() string
	GetSignerAddress() string
	GetSignerPrivateKey() *ecdsa.PrivateKey
	GetPrivateKeyHex() string
	GetChainID() int64
	NextNonce() int64
	GetRecvWindow() int64
	GetLogger() log.Logger
	GetSignFn() SignFn
	GetTimeOffsetMs() int64
	GetAuthError() error
	SetUsedWeight(used, used1m *int64)
	SetOrderCount(count10s, count1d *int64)
}

type kv struct {
	Key   string
	Value string
}

type Request struct {
	client Client
	r      *resty.Request
	method string
	path   string
	params []kv
	err    error
}

func newRequest(ctx context.Context, c Client, method, path string, params ...map[string]string) *Request {
	merged := make(map[string]string)
	for _, p := range params {
		for k, v := range p {
			merged[k] = v
		}
	}
	r := c.GetHttpClient().R().
		SetHeader("User-Agent", common.GO_ASTER_USER_AGENT).
		SetContext(ctx)
	r.Method = method
	return &Request{
		client: c,
		r:      r,
		method: method,
		path:   path,
		params: orderedKVs(merged),
	}
}

func Get(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	return newRequest(ctx, c, http.MethodGet, path, params...)
}

func Post(c Client, ctx context.Context, path string, params ...map[string]string) *Request {
	return newRequest(ctx, c, http.MethodPost, path, params...)
}

func Put(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	return newRequest(ctx, c, http.MethodPut, path, params...)
}

func Delete(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	return newRequest(ctx, c, http.MethodDelete, path, params...)
}

// WithSignature attaches the V3 EIP-712 signature to the request. Endpoints
// requiring TRADE / USER_DATA / USER_STREAM auth should chain this in their
// Do(ctx). The flow is:
//
//  1. Append nonce (microseconds, monotonic) and signer (API wallet address).
//  2. URL-encode the ordered params list -> msg.
//  3. Run msg through EIP-712 (primaryType=Message, single string field "msg").
//  4. Append the resulting signature.
//
// The final params list is later URL-encoded again in the same order to form
// the request's query string, so client-emitted bytes match the bytes the
// service used to compute its expected signature.
func (r *Request) WithSignature() *Request {
	if r.err != nil {
		return r
	}
	if err := r.ensureSigner(); err != nil {
		r.err = err
		return r
	}
	return r.signWithNonce(r.client.NextNonce())
}

// WithSignatureNonce is like WithSignature but signs over a caller-supplied
// nonce instead of a fresh monotonic one. The guarded cancel endpoints
// (guardedCancelOrder / guardedBatchOrders) require this: their nonce must
// replay the exact value submitted when the order was placed, which the engine
// uses to guard against duplicate or out-of-order cancellations.
func (r *Request) WithSignatureNonce(nonce int64) *Request {
	if r.err != nil {
		return r
	}
	if err := r.ensureSigner(); err != nil {
		r.err = err
		return r
	}
	return r.signWithNonce(nonce)
}

// ensureSigner validates that the client is ready to produce a V3 signature.
func (r *Request) ensureSigner() error {
	if err := r.client.GetAuthError(); err != nil {
		return fmt.Errorf("auth: %w", err)
	}
	if r.client.GetSignerAddress() == "" {
		return errors.New("signer not configured: call WithAuth(userAddress, signerPrivateKeyHex)")
	}
	if r.client.GetUserAddress() == "" {
		return errors.New("user address not configured: call WithAuth(userAddress, signerPrivateKeyHex)")
	}
	return nil
}

// signWithNonce appends user/nonce/signer, EIP-712-signs the ordered params,
// and appends the resulting signature.
func (r *Request) signWithNonce(nonce int64) *Request {
	signer := r.client.GetSignerAddress()
	r.params = append(r.params,
		kv{Key: "user", Value: r.client.GetUserAddress()},
		kv{Key: "nonce", Value: strconv.FormatInt(nonce, 10)},
		kv{Key: "signer", Value: signer},
	)

	msg := encodeOrdered(r.params)

	chainID := r.client.GetChainID()
	var (
		signature string
		err       error
	)
	if fn := r.client.GetSignFn(); fn != nil {
		signature, err = fn(r.client.GetPrivateKeyHex(), msg, chainID)
	} else if priv := r.client.GetSignerPrivateKey(); priv != nil {
		signature, err = signWithKey(priv, msg, chainID)
	} else {
		err = errors.New("no signing key or sign function configured")
	}
	if err != nil {
		r.err = err
		return r
	}

	r.params = append(r.params, kv{Key: "signature", Value: signature})
	return r
}

// orderedKVs converts a map into a deterministic slice. Sorting by key
// alphabetically gives a stable, easily-debuggable order. The server treats
// the request's literal query string as the EIP-712 message, so any stable
// client-side order works as long as the URL bytes and the signed bytes
// agree -- which is guaranteed here because both come from the same slice.
func orderedKVs(m map[string]string) []kv {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]kv, 0, len(keys))
	for _, k := range keys {
		out = append(out, kv{Key: k, Value: m[k]})
	}
	return out
}

// encodeOrdered URL-encodes the params in slice order. We don't use
// url.Values.Encode because it sorts keys alphabetically internally -- we
// already control the order and want to preserve it byte-for-byte.
func encodeOrdered(params []kv) string {
	if len(params) == 0 {
		return ""
	}
	var b strings.Builder
	for i, p := range params {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(url.QueryEscape(p.Key))
		b.WriteByte('=')
		b.WriteString(url.QueryEscape(p.Value))
	}
	return b.String()
}

// fullURL builds the request URL with the (possibly signed) ordered params
// appended as a query string. We assemble it ourselves rather than relying
// on resty's QueryParam (a map) so the byte order matches the signed msg.
func (r *Request) fullURL() string {
	base := strings.TrimSuffix(r.client.GetHttpClient().BaseURL, "/")
	path := r.path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	urlStr := base + path
	if q := encodeOrdered(r.params); q != "" {
		urlStr += "?" + q
	}
	return urlStr
}
