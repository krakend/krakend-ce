package payloads

import (
	"encoding/json"
	"net/http"
)

func CollectionEndpoint(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	var res []interface{}

	for i := 0; i < 10; i++ {
		res = append(res, map[string]interface{}{
			"path": r.URL.Path,
			"i":    i,
		})
	}

	json.NewEncoder(rw).Encode(res)
}
