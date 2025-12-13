package request

import (
	"fmt"

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

	// TODO response header handler
	if response.IsError() {
		// TODO: handler error
		return nil, fmt.Errorf("status code: %d, body: %s", response.StatusCode(), response.String())
	}
	return resp, nil
}
