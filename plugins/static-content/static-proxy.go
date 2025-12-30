package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
)

func handleStaticContent(config PluginConfig, w http.ResponseWriter, req *http.Request, h http.Handler) {
	logger.Debug(fmt.Sprintf("Handling static content request for %s", req.URL.Path))
	if matchesGatewayEndpoints(config.ServiceGateway, req.URL.Path) {
		h.ServeHTTP(w, req) //Continue with the gateway chain
		return
	}

	if matchedStatic := matchesStaticContentEndpoint(config.Static, req.URL.Path); matchedStatic.ServiceHost != "" {
		proxyToStaticServer(matchedStatic, w, req)
		return
	}

	h.ServeHTTP(w, req)
}

func proxyToStaticServer(staticServer StaticConfig, w http.ResponseWriter, req *http.Request) {
	targetURL, err := url.Parse(staticServer.ServiceHost)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse service host: %s", err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	proxy := &httputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			outReq.URL.Scheme = targetURL.Scheme
			outReq.URL.Host = targetURL.Host
			outReq.Host = targetURL.Host
			outReq.URL.Path = req.URL.Path
			outReq.URL.RawQuery = req.URL.RawQuery

			// Copy original headers to preserve those set by upstream proxy
			for key, values := range req.Header {
				outReq.Header[key] = values
			}

			// Append client IP to X-Forwarded-For chain
			clientIP := extractClientIP(req)
			if clientIP != "" {
				appendForwardedFor(outReq, clientIP)
			}

			// Only set X-Forwarded-Host if not already present
			if outReq.Header.Get("X-Forwarded-Host") == "" {
				outReq.Header.Set("X-Forwarded-Host", req.Host)
			}

			// Only set X-Forwarded-Proto if not already present
			if outReq.Header.Get("X-Forwarded-Proto") == "" {
				outReq.Header.Set("X-Forwarded-Proto", schemeFromRequest(req))
			}

			// Set X-Real-IP to the original client IP if not already set by Traefik
			if outReq.Header.Get("X-Real-IP") == "" {
				if originalIP := req.Header.Get("X-Forwarded-For"); originalIP != "" {
					// Use the first IP in the chain (original client)
					outReq.Header.Set("X-Real-IP", strings.TrimSpace(strings.Split(originalIP, ",")[0]))
				} else if clientIP != "" {
					outReq.Header.Set("X-Real-IP", clientIP)
				}
			}
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error(fmt.Sprintf("proxy error: %s", err.Error()))
			w.WriteHeader(http.StatusBadGateway)
		},
	}

	proxy.ServeHTTP(w, req)
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

// appendForwardedFor appends clientIP to the X-Forwarded-For header chain
func appendForwardedFor(req *http.Request, clientIP string) {
	if prior := req.Header.Get("X-Forwarded-For"); prior != "" {
		req.Header.Set("X-Forwarded-For", prior+", "+clientIP)
	} else {
		req.Header.Set("X-Forwarded-For", clientIP)
	}
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

// ... existing code ...

func matchesStaticContentEndpoint(config []StaticConfig, path string) StaticConfig {
	for _, s := range config {
		if matchesWildcard(path, s.PathPrefix) {
			return s
		}
	}
	return StaticConfig{}
}

func matchesGatewayEndpoints(config ServiceGatewayConfig, path string) bool {
	if existsInPaths(path, config.PathPrefix) {
		return true // matches at least one gateway endpoint
	}
	return false
}

func existsInPaths(path string, paths []string) bool {
	return slices.ContainsFunc(paths, func(s string) bool {
		if matchesWildcard(path, s) {
			return true
		}
		return false
	})
}

func matchesWildcard(s, wildcard string) bool {
	return strings.HasPrefix(s, strings.TrimSuffix(wildcard, "*"))
}
