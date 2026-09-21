package krakend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	serverhttp "github.com/luraproject/lura/v2/transport/http/server"
)

// connectionLeaseNamespace is the service-level extra_config key that configures
// the leasing strategy of the shared outbound (backend) HTTP connection pool.
//
//	"extra_config": {
//	  "backend/http/client": { "connection_lease_strategy": "fifo" }
//	}
const connectionLeaseNamespace = "backend/http/client"

const (
	// leaseStrategyLIFO is the Go net/http default: the most-recently-used idle
	// connection is reused first (stack). It is the behaviour when the strategy
	// is unset, so no transport is replaced for it.
	leaseStrategyLIFO = "lifo"
	// leaseStrategyFIFO distributes requests across a ring of independent
	// transports so idle connections are cycled evenly and traffic is spread
	// across backend nodes (queue-like behaviour). Go's stdlib transport cannot
	// do FIFO within a single idle pool, so we round-robin over N pools to obtain
	// the even-distribution benefits (Option A).
	leaseStrategyFIFO = "fifo"

	// defaultConnectionRingSize is the number of independent transports in the
	// FIFO ring when "connection_pools" is not set. Eight suits a handful of
	// backend replicas while keeping strong keep-alive reuse under the default
	// max_idle_connections_per_host (250).
	defaultConnectionRingSize = 8
	// maxConnectionRingSize is an upper safety clamp so a mistyped config cannot
	// spawn an unreasonable number of transports/sockets.
	maxConnectionRingSize = 256

	connectionLeaseLogPrefix = "[SERVICE: Backend Connection Pool]"
)

// connectionLeaseConfig is the parsed service-level configuration for the
// outbound connection leasing strategy.
type connectionLeaseConfig struct {
	// Strategy selects the leasing behaviour: "lifo" (default) or "fifo".
	Strategy string `json:"connection_lease_strategy"`
	// Pools is the number of independent transports (idle-connection pools) the
	// "fifo" strategy round-robins across. This is the distribution fan-out, NOT
	// the idle-connection count: the existing max_idle_connections* limits are
	// split across these pools. Rule of thumb: set it around the number of
	// backend replicas behind the load balancer, keeping it <=
	// max_idle_connections_per_host / 2. Ignored for "lifo"; defaults to
	// defaultConnectionRingSize.
	Pools int `json:"connection_pools"`
}

// connectionLeaseConfigGetter extracts the connection lease config from the
// service extra_config. The second return value reports whether the namespace
// was present at all.
func connectionLeaseConfigGetter(e config.ExtraConfig) (connectionLeaseConfig, bool) {
	v, ok := e[connectionLeaseNamespace]
	if !ok {
		return connectionLeaseConfig{}, false
	}
	b, err := json.Marshal(v)
	if err != nil {
		return connectionLeaseConfig{}, false
	}
	cfg := connectionLeaseConfig{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return connectionLeaseConfig{}, false
	}
	return cfg, true
}

// setupBackendConnectionLease configures the leasing strategy of the shared
// outbound HTTP connection pool used by every backend.
//
// The default (and "lifo") behaviour is a no-op: the router later initialises
// http.DefaultTransport with lura's stdlib transport (LIFO idle reuse).
//
// For the "fifo" strategy we build lura's transport from the service config, use
// it as a template for a ring of independent clones, and install a round-robin
// RoundTripper as http.DefaultTransport. Building the default transport here
// consumes lura's one-shot initialisation, so the router's later
// InitHTTPDefaultTransport call becomes a no-op and our ring survives.
//
// This must be called before the router is started (i.e. before
// routerFactory...Run) while the service config is still in scope.
func setupBackendConnectionLease(cfg config.ServiceConfig, logger logging.Logger) {
	leaseCfg, ok := connectionLeaseConfigGetter(cfg.ExtraConfig)
	if !ok {
		return
	}

	switch strings.ToLower(strings.TrimSpace(leaseCfg.Strategy)) {
	case "", leaseStrategyLIFO:
		// stdlib default; let the router configure the transport as usual.
		return
	case leaseStrategyFIFO:
		// handled below
	default:
		logger.Warning(fmt.Sprintf("%s unknown connection_lease_strategy %q, falling back to lifo",
			connectionLeaseLogPrefix, leaseCfg.Strategy))
		return
	}

	size := leaseCfg.Pools
	if size <= 1 {
		size = defaultConnectionRingSize
	}
	if size > maxConnectionRingSize {
		logger.Warning(fmt.Sprintf("%s connection_pools %d exceeds the maximum %d, clamping",
			connectionLeaseLogPrefix, size, maxConnectionRingSize))
		size = maxConnectionRingSize
	}

	// Build and freeze the default transport from the service config, then read
	// it back as the ring template. This also consumes lura's onceTransportConfig
	// so the router's later InitHTTPDefaultTransport call is a no-op.
	serverhttp.InitHTTPDefaultTransportWithLogger(cfg, logger)

	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		logger.Warning(fmt.Sprintf("%s default transport is not *http.Transport; keeping lifo",
			connectionLeaseLogPrefix))
		return
	}

	// A ring wider than the per-host idle cap starves each pool of keep-alive
	// reuse (every pool holds a single idle connection and re-dials per request).
	if base.MaxIdleConnsPerHost > 0 && size > base.MaxIdleConnsPerHost {
		logger.Warning(fmt.Sprintf("%s connection_pools %d exceeds max_idle_connections_per_host %d; each pool keeps a single idle connection, reducing keep-alive reuse",
			connectionLeaseLogPrefix, size, base.MaxIdleConnsPerHost))
	}

	ring := make([]*http.Transport, size)
	for i := range ring {
		clone := base.Clone()
		// Keep the aggregate idle-pool capacity roughly aligned with the
		// configured limits by splitting them across the ring. A value of 0
		// means "unlimited" in net/http, so it is left untouched.
		if base.MaxIdleConns > 0 {
			clone.MaxIdleConns = ceilDiv(base.MaxIdleConns, size)
		}
		if base.MaxIdleConnsPerHost > 0 {
			clone.MaxIdleConnsPerHost = ceilDiv(base.MaxIdleConnsPerHost, size)
		}
		ring[i] = clone
	}

	http.DefaultTransport = newRoundRobinTransport(ring)

	logger.Info(fmt.Sprintf("%s FIFO connection leasing enabled (ring of %d transports)",
		connectionLeaseLogPrefix, size))
}

// ceilDiv returns ceil(a/b) with a floor of 1, guarding against a zero divisor.
func ceilDiv(a, b int) int {
	if b <= 0 {
		return a
	}
	r := a / b
	if a%b != 0 {
		r++
	}
	if r < 1 {
		r = 1
	}
	return r
}

// roundRobinTransport dispatches each request to the next transport in the ring
// using an atomic counter. Each transport keeps its own idle-connection pool, so
// spreading requests across the ring cycles connections evenly and distributes
// traffic across backend nodes instead of always reusing the warmest connection.
type roundRobinTransport struct {
	transports []*http.Transport
	counter    uint64
}

func newRoundRobinTransport(transports []*http.Transport) *roundRobinTransport {
	return &roundRobinTransport{transports: transports}
}

// pick returns the next transport in the ring.
func (rr *roundRobinTransport) pick() *http.Transport {
	n := atomic.AddUint64(&rr.counter, 1) - 1
	return rr.transports[int(n%uint64(len(rr.transports)))]
}

// RoundTrip implements http.RoundTripper.
func (rr *roundRobinTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return rr.pick().RoundTrip(req)
}

// CloseIdleConnections closes idle connections on every transport in the ring so
// the ring honours the http.Client.CloseIdleConnections contract.
func (rr *roundRobinTransport) CloseIdleConnections() {
	for _, t := range rr.transports {
		t.CloseIdleConnections()
	}
}
