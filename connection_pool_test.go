package krakend

import (
	"net/http"
	"testing"

	"github.com/luraproject/lura/v2/config"
)

func TestConnectionLeaseConfigGetter(t *testing.T) {
	for _, tc := range []struct {
		name      string
		extra     config.ExtraConfig
		wantOK    bool
		wantStrat string
		wantPools int
	}{
		{name: "absent", extra: config.ExtraConfig{}, wantOK: false},
		{
			name:      "fifo with pools",
			extra:     config.ExtraConfig{connectionLeaseNamespace: map[string]interface{}{"connection_lease_strategy": "fifo", "connection_pools": 4}},
			wantOK:    true,
			wantStrat: "fifo",
			wantPools: 4,
		},
		{
			name:      "lifo",
			extra:     config.ExtraConfig{connectionLeaseNamespace: map[string]interface{}{"connection_lease_strategy": "lifo"}},
			wantOK:    true,
			wantStrat: "lifo",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := connectionLeaseConfigGetter(tc.extra)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.Strategy != tc.wantStrat {
				t.Errorf("strategy = %q, want %q", got.Strategy, tc.wantStrat)
			}
			if got.Pools != tc.wantPools {
				t.Errorf("pools = %d, want %d", got.Pools, tc.wantPools)
			}
		})
	}
}

func TestCeilDiv(t *testing.T) {
	cases := map[string][3]int{ // input a, b -> want
		"exact":    {250, 5, 50},
		"round up": {250, 8, 32},
		"floor 1":  {3, 8, 1},
		"zero b":   {250, 0, 250},
	}
	for name, c := range cases {
		if got := ceilDiv(c[0], c[1]); got != c[2] {
			t.Errorf("%s: ceilDiv(%d,%d) = %d, want %d", name, c[0], c[1], got, c[2])
		}
	}
}

// TestRoundRobinTransportRotates verifies pick() cycles through the ring's
// distinct transports in order.
func TestRoundRobinTransportRotates(t *testing.T) {
	const size = 3
	ring := make([]*http.Transport, size)
	for i := range ring {
		ring[i] = &http.Transport{}
	}
	rr := newRoundRobinTransport(ring)

	seen := make([]*http.Transport, 2*size)
	for i := range seen {
		seen[i] = rr.pick()
	}

	// round-robin: pick i and pick i+size must be the same transport.
	for i := 0; i < size; i++ {
		if seen[i] != seen[i+size] {
			t.Errorf("pick %d and %d differ; expected round-robin", i, i+size)
		}
	}
	// the transports in one cycle must be distinct pools.
	if seen[0] == seen[1] || seen[1] == seen[2] || seen[0] == seen[2] {
		t.Errorf("expected %d distinct transports in the ring", size)
	}
	if rr.counter != uint64(len(seen)) {
		t.Errorf("counter = %d, want %d", rr.counter, len(seen))
	}
}
