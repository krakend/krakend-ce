package checks

import (
	"net"
	"net/http"
)

func XForwardedFor(h http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if ip := net.ParseIP(r.Header.Get("X-Forwarded-For")); ip == nil || !ip.IsLoopback() {
			http.Error(rw, "invalid X-Forwarded-For", 400)
			return
		}
		h.ServeHTTP(rw, r)
	}
}
