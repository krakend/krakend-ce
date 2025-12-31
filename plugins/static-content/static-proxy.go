package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"headers"
	"paths"
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

			headers.CopyAll(outReq, req)
			headers.SetProxyHeaders(outReq, req)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error(fmt.Sprintf("proxy error: %s", err.Error()))
			w.WriteHeader(http.StatusBadGateway)
		},
	}

	proxy.ServeHTTP(w, req)
}

func matchesStaticContentEndpoint(config []StaticConfig, path string) StaticConfig {
	for _, s := range config {
		if paths.MatchesWildcard(path, s.PathPrefix) {
			return s
		}
	}
	return StaticConfig{}
}

func matchesGatewayEndpoints(config ServiceGatewayConfig, path string) bool {
	if paths.ExistsInPaths(path, config.PathPrefix) {
		return true // matches at least one gateway endpoint
	}
	return false
}
