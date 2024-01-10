package payloads

import (
	"encoding/json"
	"io"
	"net/http"
)

func EchoBuilder() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"path":    r.URL.Path,
			"query":   r.URL.Query(),
			"headers": r.Header,
			"foo":     42,
		}

		if r.URL.Query().Get("dump_body") == "1" {
			b, _ := io.ReadAll(r.Body)
			r.Body.Close()
			resp["body"] = string(b)
		}

		json.NewEncoder(rw).Encode(resp)
	}
}

func DefaultEchoEndpoint() http.HandlerFunc {
	return AddResponseHeaders(
		RemoveRequestHeader(
			EchoBuilder(), []string{"X-Forwarded-For"}),
		http.Header{
			"Content-Type": []string{"application/json"},
			"Set-Cookie": []string{
				"test1=test1",
				"test2=test2",
			}})
}
