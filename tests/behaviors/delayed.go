package behaviors

import (
	"net/http"
	"time"
)

func Delayed(d time.Duration, h http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		<-time.After(d)
		h.ServeHTTP(rw, req)
	}
}
