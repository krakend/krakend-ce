package tests

import (
	"fmt"
	"net/http"

	"github.com/krakendio/krakend-ce/v2/tests/behaviors"
	"github.com/krakendio/krakend-ce/v2/tests/checks"
	"github.com/krakendio/krakend-ce/v2/tests/payloads"
)

var DefaultBackendBuilder mockBackendBuilder

type mockBackendBuilder struct{}

func (mockBackendBuilder) New(cfg *Config) http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/param_forwarding/", checks.XForwardedFor(http.HandlerFunc(payloads.DefaultEchoEndpoint())))
	mux.HandleFunc("/xml", checks.XForwardedFor(http.HandlerFunc(payloads.XmlEndpoint)))
	mux.HandleFunc("/collection/", checks.XForwardedFor(http.HandlerFunc(payloads.CollectionEndpoint)))
	mux.HandleFunc("/delayed/", checks.XForwardedFor(behaviors.Delayed(cfg.getDelay(), payloads.DefaultEchoEndpoint())))
	mux.HandleFunc("/redirect/", checks.XForwardedFor(http.HandlerFunc(behaviors.RedirectEndpoint)))
	mux.HandleFunc("/jwk/symmetric", http.HandlerFunc(payloads.SymmetricJWKEndpoint))

	return http.Server{ // skipcq: GO-S2112
		Addr:    fmt.Sprintf(":%v", cfg.getBackendPort()),
		Handler: mux,
	}
}
