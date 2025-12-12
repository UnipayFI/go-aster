package request

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/UnipayFI/go-aster/common"
	"github.com/UnipayFI/go-aster/pkg/log"
	"github.com/go-resty/resty/v2"
)

type Client interface {
	GetHttpClient() *resty.Client
	GetRecvWindow() int64
	GetApiKey() string
	GetSecretKey() string
	GetLogger() log.Logger
	GetSignFn() func(apiKey, secretKey string, payload string) (string, error)
}

type Request struct {
	client Client
	r      *resty.Request
	err    error
}

func Get(ctx context.Context, client Client, url string, params ...map[string]string) *Request {
	r := new(ctx, client, http.MethodGet, url)
	for _, param := range params {
		r.SetQueryParams(param)
	}
	return &Request{client: client, r: r}
}

func Post(client Client, ctx context.Context, url string, body any) *Request {
	r := new(ctx, client, http.MethodPost, url)
	r.SetBody(body)
	return &Request{client: client, r: r}
}

func (r *Request) Sign() *Request {
	// 1. set recvWindow and timestamp
	r.setBaseParams()

	// 2. get payload
	payload, err := r.payload()
	if err != nil {
		r.err = err
		return r
	}

	// 3. sign payload
	signFn := r.client.GetSignFn()
	if signFn == nil {
		signFn = sign
	}
	signature, err := signFn(r.client.GetApiKey(), r.client.GetSecretKey(), payload)
	if err != nil {
		r.err = err
		return r
	}

	// 4. set signature
	if r.r.Method == http.MethodGet {
		r.r.SetQueryParam("signature", signature)
	} else {
		r.r.Body.(map[string]any)["signature"] = signature
	}
	return r
}

func (r *Request) setBaseParams() {
	r.r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	r.r.SetHeader("X-MBX-APIKEY", r.client.GetApiKey())

	switch r.r.Method {
	case http.MethodGet:
		r.r.SetQueryParam("recvWindow", fmt.Sprintf("%d", r.client.GetRecvWindow()))
		r.r.SetQueryParam("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	case http.MethodPost:
		body := r.r.Body.(map[string]any)
		body["recvWindow"] = fmt.Sprintf("%d", r.client.GetRecvWindow())
		body["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
		r.r.Body = body
	}
}

func (r *Request) payload() (payload string, err error) {
	switch r.r.Method {
	case http.MethodGet:
		payload = r.r.QueryParam.Encode()
		return payload, nil
	case http.MethodPost:
		body, err := r.client.GetHttpClient().JSONMarshal(r.r.Body)
		if err != nil {
			return "", err
		}
		payload = string(body)
		return payload, nil
	default:
		return "", fmt.Errorf("invalid method: %s", r.r.Method)
	}
}

func new(ctx context.Context, client Client, method, url string) *resty.Request {
	r := client.GetHttpClient().R().
		SetHeader("User-Agent", common.GO_ASTER_USER_AGENT).
		SetContext(ctx)
	r.Method = method
	r.URL = url
	return r
}

func sign(apiKey, secretKey string, payload string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", (mac.Sum(nil))), nil
}
