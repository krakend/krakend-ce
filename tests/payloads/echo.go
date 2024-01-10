package payloads

import (
	"encoding/json"
	"io"
	"net/http"
)

func EchoEndpoint(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	rw.Header().Add("Set-Cookie", "test1=test1")
	rw.Header().Add("Set-Cookie", "test2=test2")
	r.Header.Del("X-Forwarded-For")
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
