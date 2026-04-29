package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnipayFI/go-aster/v3/client"
	"github.com/go-resty/resty/v2"
)

// Do executes the request and decodes the response body into *T. Following
// the official V3 Python demo (send_by_url), all parameters -- including the
// signature -- are sent as URL query string for every method (GET/POST/PUT/
// DELETE); the body is left empty.
func Do[T any](r *Request) (resp *T, err error) {
	if r.err != nil {
		return nil, r.err
	}

	r.r.URL = r.fullURL()
	r.client.GetLogger().Debugf("request: %s %s", r.method, r.r.URL)
	defer func() {
		if err == nil {
			return
		}
		r.client.GetLogger().Errorf("request %s failed: %s", r.r.URL, err)
	}()

	var response *resty.Response
	response, err = r.r.SetResult(&resp).Send()
	if err != nil {
		return nil, err
	}

	r.client.GetLogger().Debugf("response: %v", response.String())
	defer response.RawBody().Close()

	handlerRateLimit(r, response)
	if response.IsError() {
		return nil, handlerAPIError(r, response)
	}
	return resp, nil
}

// DoRaw is like Do but returns the raw response body without JSON decoding.
// Used by endpoints that return non-JSON or non-uniform shapes (rare).
func DoRaw(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	r.r.URL = r.fullURL()
	r.client.GetLogger().Debugf("request: %s %s", r.method, r.r.URL)

	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	defer response.RawBody().Close()

	handlerRateLimit(r, response)
	if response.IsError() {
		return nil, handlerAPIError(r, response)
	}
	return response.Body(), nil
}

func handlerRateLimit(r *Request, response *resty.Response) {
	handler := func(header http.Header, key string) *int64 {
		if value := header.Get(key); value != "" {
			if used, err := strconv.ParseInt(value, 10, 64); err == nil {
				return &used
			}
		}
		return nil
	}
	used := handler(response.Header(), "X-Mbx-Used-Weight")
	used1m := handler(response.Header(), "X-Mbx-Used-Weight-1m")
	count10s := handler(response.Header(), "X-Mbx-Order-Count-10s")
	count1d := handler(response.Header(), "X-Mbx-Order-Count-1d")
	r.client.SetUsedWeight(used, used1m)
	r.client.SetOrderCount(count10s, count1d)
}

func handlerAPIError(r *Request, response *resty.Response) error {
	apiErr := &client.APIError{}
	e := json.Unmarshal(response.Body(), apiErr)
	if e != nil {
		r.client.GetLogger().Errorf("failed to unmarshal json: %s\n", e)
	}
	if !apiErr.IsValid() {
		return fmt.Errorf("request failed with status code: %d, body: %s", response.StatusCode(), response.String())
	}
	return apiErr
}
