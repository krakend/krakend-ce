package behaviors

import (
	"net/http"
	"strconv"
)

func RedirectEndpoint(rw http.ResponseWriter, r *http.Request) {
	u := r.URL
	u.Path = "/param_forwarding/"

	status, ok2 := r.URL.Query()["status"]
	code := 301
	if !ok2 || status[0] != "301" {
		var err error
		code, err = strconv.Atoi(status[0])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	http.Redirect(rw, r, u.String(), code)
}
