package headers

import (
	"net"
	"net/http"
	"strings"
)

func CopyAll(outReq *http.Request, req *http.Request) {
	for key, values := range req.Header {
		outReq.Header[key] = values
	}
}

func SetProxyHeaders(outReq *http.Request, req *http.Request) {
	clientIP := proxyForward(outReq, req)

	// Set X-Real-IP to the original client IP if not already set
	if outReq.Header.Get("X-Real-IP") == "" {
		if originalIP := req.Header.Get("X-Forwarded-For"); originalIP != "" {
			// Use the first IP in the chain (original client)
			outReq.Header.Set("X-Real-IP", strings.TrimSpace(strings.Split(originalIP, ",")[0]))
		} else if clientIP != "" {
			outReq.Header.Set("X-Real-IP", clientIP)
		}
	}
}

func proxyForward(outReq *http.Request, req *http.Request) string {
	// Start with existing X-Forwarded-For from incoming request
	if incomingXFF := req.Header.Get("X-Forwarded-For"); incomingXFF != "" {
		outReq.Header.Set("X-Forwarded-For", incomingXFF)
	}

	// Append current client IP to the chain
	clientIP := extractClientIP(req)
	if clientIP != "" {
		appendForwardedFor(outReq, clientIP)
	}

	proxyPassXForwarded(outReq, req)
	return clientIP
}

// AppendForwardedFor appends clientIP to the X-Forwarded-For header chain
func appendForwardedFor(req *http.Request, clientIP string) {
	if prior := req.Header.Get("X-Forwarded-For"); prior != "" {
		req.Header.Set("X-Forwarded-For", prior+", "+clientIP)
	} else {
		req.Header.Set("X-Forwarded-For", clientIP)
	}
}

func proxyPassXForwarded(outReq *http.Request, req *http.Request) {
	// Only set X-Forwarded-Host if not already present
	if outReq.Header.Get("X-Forwarded-Host") == "" {
		outReq.Header.Set("X-Forwarded-Host", req.Host)
	}

	// Only set X-Forwarded-Proto if not already present
	if outReq.Header.Get("X-Forwarded-Proto") == "" {
		outReq.Header.Set("X-Forwarded-Proto", schemeFromRequest(req))
	}
}

// extractClientIP extracts the IP address from RemoteAddr, stripping the port
func extractClientIP(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		// RemoteAddr might not have a port
		return req.RemoteAddr
	}
	return host
}

func schemeFromRequest(req *http.Request) string {
	// Respect existing X-Forwarded-Proto from upstream proxy
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if req.TLS != nil {
		return "https"
	}
	return "http"
}
