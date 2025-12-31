package headers

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyAll(t *testing.T) {
	tests := []struct {
		name            string
		sourceHeaders   map[string][]string
		existingHeaders map[string][]string
	}{
		{
			name: "copy single header",
			sourceHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			existingHeaders: map[string][]string{},
		},
		{
			name: "copy multiple headers",
			sourceHeaders: map[string][]string{
				"Content-Type":    {"application/json"},
				"Authorization":   {"Bearer token"},
				"X-Custom-Header": {"custom-value"},
			},
			existingHeaders: map[string][]string{},
		},
		{
			name: "copy headers with multiple values",
			sourceHeaders: map[string][]string{
				"Accept": {"text/html", "application/json"},
			},
			existingHeaders: map[string][]string{},
		},
		{
			name: "overwrite existing headers",
			sourceHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			existingHeaders: map[string][]string{
				"Content-Type": {"text/plain"},
			},
		},
		{
			name:          "empty source headers",
			sourceHeaders: map[string][]string{},
			existingHeaders: map[string][]string{
				"Existing": {"value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Header: make(http.Header)}
			for k, v := range tt.sourceHeaders {
				req.Header[k] = v
			}

			outReq := &http.Request{Header: make(http.Header)}
			for k, v := range tt.existingHeaders {
				outReq.Header[k] = v
			}

			CopyAll(outReq, req)

			for k, expectedValues := range tt.sourceHeaders {
				actualValues := outReq.Header[k]
				assert.Equal(t, expectedValues, actualValues, "header %s mismatch", k)
			}
		})
	}
}

func TestSetProxyHeaders(t *testing.T) {
	tests := []struct {
		name                 string
		remoteAddr           string
		incomingHeaders      map[string]string
		existingOutHeaders   map[string]string
		expectedForwardedFor string
		expectedRealIP       string
		expectedHost         string
		expectedProto        string
		useTLS               bool
	}{
		{
			name:                 "basic proxy with IP and port",
			remoteAddr:           "192.168.1.100:54321",
			incomingHeaders:      map[string]string{},
			existingOutHeaders:   map[string]string{},
			expectedForwardedFor: "192.168.1.100",
			expectedRealIP:       "192.168.1.100",
			expectedHost:         "example.com",
			expectedProto:        "http",
		},
		{
			name:       "chain X-Forwarded-For and set X-Real-IP",
			remoteAddr: "10.0.0.5:12345",
			incomingHeaders: map[string]string{
				"X-Forwarded-For": "1.2.3.4, 5.6.7.8",
			},
			existingOutHeaders:   map[string]string{},
			expectedForwardedFor: "1.2.3.4, 5.6.7.8, 10.0.0.5",
			expectedRealIP:       "1.2.3.4",
			expectedHost:         "example.com",
			expectedProto:        "http",
		},
		{
			name:       "preserve existing X-Real-IP",
			remoteAddr: "10.0.0.5:12345",
			incomingHeaders: map[string]string{
				"X-Forwarded-For": "1.2.3.4",
			},
			existingOutHeaders: map[string]string{
				"X-Real-IP": "preserved.ip",
			},
			expectedForwardedFor: "1.2.3.4, 10.0.0.5",
			expectedRealIP:       "preserved.ip",
			expectedHost:         "example.com",
			expectedProto:        "http",
		},
		{
			name:            "preserve existing X-Forwarded-Host",
			remoteAddr:      "10.0.0.5:12345",
			incomingHeaders: map[string]string{},
			existingOutHeaders: map[string]string{
				"X-Forwarded-Host": "original.host",
			},
			expectedForwardedFor: "10.0.0.5",
			expectedRealIP:       "10.0.0.5",
			expectedHost:         "original.host",
			expectedProto:        "http",
		},
		{
			name:            "preserve existing X-Forwarded-Proto",
			remoteAddr:      "10.0.0.5:12345",
			incomingHeaders: map[string]string{},
			existingOutHeaders: map[string]string{
				"X-Forwarded-Proto": "https",
			},
			expectedForwardedFor: "10.0.0.5",
			expectedRealIP:       "10.0.0.5",
			expectedHost:         "example.com",
			expectedProto:        "https",
		},
		{
			name:                 "TLS request sets https proto",
			remoteAddr:           "10.0.0.5:12345",
			incomingHeaders:      map[string]string{},
			existingOutHeaders:   map[string]string{},
			useTLS:               true,
			expectedForwardedFor: "10.0.0.5",
			expectedRealIP:       "10.0.0.5",
			expectedHost:         "example.com",
			expectedProto:        "https",
		},
		{
			name:       "respect upstream X-Forwarded-Proto",
			remoteAddr: "10.0.0.5:12345",
			incomingHeaders: map[string]string{
				"X-Forwarded-Proto": "https",
			},
			existingOutHeaders:   map[string]string{},
			expectedForwardedFor: "10.0.0.5",
			expectedRealIP:       "10.0.0.5",
			expectedHost:         "example.com",
			expectedProto:        "https",
		},
		{
			name:                 "IPv6 address",
			remoteAddr:           "[2001:db8::1]:8080",
			incomingHeaders:      map[string]string{},
			existingOutHeaders:   map[string]string{},
			expectedForwardedFor: "2001:db8::1",
			expectedRealIP:       "2001:db8::1",
			expectedHost:         "example.com",
			expectedProto:        "http",
		},
		{
			name:                 "RemoteAddr without port",
			remoteAddr:           "192.168.1.100",
			incomingHeaders:      map[string]string{},
			existingOutHeaders:   map[string]string{},
			expectedForwardedFor: "192.168.1.100",
			expectedRealIP:       "192.168.1.100",
			expectedHost:         "example.com",
			expectedProto:        "http",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/path", nil)
			req.RemoteAddr = tt.remoteAddr
			req.Host = "example.com"

			if tt.useTLS {
				req.TLS = &tls.ConnectionState{}
			}

			for k, v := range tt.incomingHeaders {
				req.Header.Set(k, v)
			}

			outReq := httptest.NewRequest("GET", "http://backend.local/path", nil)
			for k, v := range tt.existingOutHeaders {
				outReq.Header.Set(k, v)
			}

			SetProxyHeaders(outReq, req)

			assert.Equal(t, tt.expectedForwardedFor, outReq.Header.Get("X-Forwarded-For"), "X-Forwarded-For mismatch")
			assert.Equal(t, tt.expectedRealIP, outReq.Header.Get("X-Real-IP"), "X-Real-IP mismatch")
			assert.Equal(t, tt.expectedHost, outReq.Header.Get("X-Forwarded-Host"), "X-Forwarded-Host mismatch")
			assert.Equal(t, tt.expectedProto, outReq.Header.Get("X-Forwarded-Proto"), "X-Forwarded-Proto mismatch")
		})
	}
}

func TestAppendForwardedFor(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		clientIP string
		expected string
	}{
		{
			name:     "first IP in chain",
			existing: "",
			clientIP: "192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "append to existing chain",
			existing: "1.2.3.4",
			clientIP: "5.6.7.8",
			expected: "1.2.3.4, 5.6.7.8",
		},
		{
			name:     "append to long chain",
			existing: "1.2.3.4, 5.6.7.8, 9.10.11.12",
			clientIP: "13.14.15.16",
			expected: "1.2.3.4, 5.6.7.8, 9.10.11.12, 13.14.15.16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Header: make(http.Header)}
			if tt.existing != "" {
				req.Header.Set("X-Forwarded-For", tt.existing)
			}

			appendForwardedFor(req, tt.clientIP)

			assert.Equal(t, tt.expected, req.Header.Get("X-Forwarded-For"))
		})
	}
}

func TestExtractClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		expected   string
	}{
		{
			name:       "IPv4 with port",
			remoteAddr: "192.168.1.100:54321",
			expected:   "192.168.1.100",
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[2001:db8::1]:8080",
			expected:   "2001:db8::1",
		},
		{
			name:       "IP without port",
			remoteAddr: "192.168.1.100",
			expected:   "192.168.1.100",
		},
		{
			name:       "localhost with port",
			remoteAddr: "127.0.0.1:12345",
			expected:   "127.0.0.1",
		},
		{
			name:       "IPv6 localhost with port",
			remoteAddr: "[::1]:8080",
			expected:   "::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{RemoteAddr: tt.remoteAddr}
			actual := extractClientIP(req)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSchemeFromRequest(t *testing.T) {
	tests := []struct {
		name      string
		useTLS    bool
		xfpHeader string
		expected  string
	}{
		{
			name:     "http without TLS",
			useTLS:   false,
			expected: "http",
		},
		{
			name:     "https with TLS",
			useTLS:   true,
			expected: "https",
		},
		{
			name:      "respect upstream X-Forwarded-Proto https",
			useTLS:    false,
			xfpHeader: "https",
			expected:  "https",
		},
		{
			name:      "respect upstream X-Forwarded-Proto http",
			useTLS:    true,
			xfpHeader: "http",
			expected:  "http",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)

			if tt.useTLS {
				req.TLS = &tls.ConnectionState{}
			}

			if tt.xfpHeader != "" {
				req.Header.Set("X-Forwarded-Proto", tt.xfpHeader)
			}

			actual := schemeFromRequest(req)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestProxyForward(t *testing.T) {
	tests := []struct {
		name               string
		remoteAddr         string
		incomingXFF        string
		expectedXFF        string
		expectedReturnedIP string
	}{
		{
			name:               "basic forward",
			remoteAddr:         "192.168.1.1:1234",
			incomingXFF:        "",
			expectedXFF:        "192.168.1.1",
			expectedReturnedIP: "192.168.1.1",
		},
		{
			name:               "chain forwarding",
			remoteAddr:         "10.0.0.1:5678",
			incomingXFF:        "1.2.3.4",
			expectedXFF:        "1.2.3.4, 10.0.0.1",
			expectedReturnedIP: "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = tt.remoteAddr
			req.Host = "example.com"

			if tt.incomingXFF != "" {
				req.Header.Set("X-Forwarded-For", tt.incomingXFF)
			}

			outReq := httptest.NewRequest("GET", "http://backend.local", nil)

			returnedIP := proxyForward(outReq, req)

			assert.Equal(t, tt.expectedReturnedIP, returnedIP)
			assert.Equal(t, tt.expectedXFF, outReq.Header.Get("X-Forwarded-For"))
			assert.Equal(t, "example.com", outReq.Header.Get("X-Forwarded-Host"))
		})
	}
}

func TestProxyPassXForwarded(t *testing.T) {
	tests := []struct {
		name             string
		reqHost          string
		reqTLS           bool
		incomingXFP      string
		existingOutHost  string
		existingOutProto string
		expectedHost     string
		expectedProto    string
	}{
		{
			name:          "set both headers",
			reqHost:       "example.com",
			reqTLS:        false,
			expectedHost:  "example.com",
			expectedProto: "http",
		},
		{
			name:          "https from TLS",
			reqHost:       "secure.com",
			reqTLS:        true,
			expectedHost:  "secure.com",
			expectedProto: "https",
		},
		{
			name:            "preserve existing host",
			reqHost:         "new.com",
			existingOutHost: "existing.com",
			expectedHost:    "existing.com",
			expectedProto:   "http",
		},
		{
			name:             "preserve existing proto",
			reqHost:          "example.com",
			existingOutProto: "https",
			expectedHost:     "example.com",
			expectedProto:    "https",
		},
		{
			name:          "respect upstream XFP",
			reqHost:       "example.com",
			reqTLS:        false,
			incomingXFP:   "https",
			expectedHost:  "example.com",
			expectedProto: "https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+tt.reqHost, nil)
			req.Host = tt.reqHost

			if tt.reqTLS {
				req.TLS = &tls.ConnectionState{}
			}

			if tt.incomingXFP != "" {
				req.Header.Set("X-Forwarded-Proto", tt.incomingXFP)
			}

			outReq := httptest.NewRequest("GET", "http://backend.local", nil)

			if tt.existingOutHost != "" {
				outReq.Header.Set("X-Forwarded-Host", tt.existingOutHost)
			}

			if tt.existingOutProto != "" {
				outReq.Header.Set("X-Forwarded-Proto", tt.existingOutProto)
			}

			proxyPassXForwarded(outReq, req)

			assert.Equal(t, tt.expectedHost, outReq.Header.Get("X-Forwarded-Host"))
			assert.Equal(t, tt.expectedProto, outReq.Header.Get("X-Forwarded-Proto"))
		})
	}
}
