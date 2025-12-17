package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnipayFI/go-aster/client"
	"github.com/go-resty/resty/v2"
)

func Do[T any](r *Request) (resp *T, err error) {
	if r.err != nil {
		return nil, r.err
	}

	r.client.GetLogger().Debugf("request: %s %s", r.r.Method, r.r.URL)
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
