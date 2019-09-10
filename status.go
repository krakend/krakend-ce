package krakend

import (
	"fmt"
	"bytes"
	"context"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/transport/http/client"
)

func GetHTTPStatusHandler(cfg *config.Backend) client.HTTPStatusHandler {
	if e, ok := cfg.ExtraConfig[client.Namespace]; ok {
		if m, ok := e.(map[string]interface{}); ok {
			if v, ok := m["error_namespace"]; ok {
				if name, ok := v.(string); ok && name != "" {
					return ErrorHTTPStatusHandler(client.DefaultHTTPStatusHandler, name)
				}
			}
		}
	}
	return client.GetHTTPStatusHandler(cfg)
}

func ErrorHTTPStatusHandler(next client.HTTPStatusHandler, name string) client.HTTPStatusHandler {
	return func(ctx context.Context, resp *http.Response) (*http.Response, error) {
		if r, err := next(ctx, resp); err == nil {
			return r, nil
		}

		var body interface{}
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			if err := json.Unmarshal(data, &body); err != nil {
				body = string(data)
			}
		}

		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		return resp, &ResponseError{
			Code: resp.StatusCode,
			Body: body,
			name: name,
		}
	}
}

type ResponseError struct {
	Code int `json:"status"`
	Body interface{} `json:"body,omitempty"`
	name string
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("request failed with status %d", r.Code)
}

func (r *ResponseError) Name() string {
	return r.name
}

func (r *ResponseError) StatusCode() int {
	return r.Code
}
