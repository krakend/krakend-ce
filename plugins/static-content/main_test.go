package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validExtraConfig = "{" +
	"\"extra_config\": " +
	"{" +
	"\"plugin/http-server\":" +
	"{" +
	"\"name\":[\"hog-static-content\"]," +
	"\"hog-static-content\":" +
	"{" +
	"\"static\":" +
	"[" +
	"{" +
	"\"path-prefix\": \"/*\"," +
	"\"service-host\": \"http://web-example\"" +
	"}" +
	"]," +
	"\"service-gateway\": " +
	"{" +
	"\"path-prefix\": [\"/api/*\"]" +
	"}" +
	"}" +
	"}" +
	"}" +
	"}"

func TestCanParseExtraPluginConfig(t *testing.T) {
	var cfg map[string]interface{}
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	if err != nil {
		t.Error(err)
	}

	var loadedPluginConfig PluginConfig
	loadedPluginConfig, err = loadPluginConfig(cfg)

	assert.NoError(t, err)
	assert.Len(t, loadedPluginConfig.Static, 1)
	assert.Equal(t, "/*", loadedPluginConfig.Static[0].PathPrefix)
	assert.Equal(t, "http://web-example", loadedPluginConfig.Static[0].ServiceHost)
	assert.Len(t, loadedPluginConfig.ServiceGateway.PathPrefix, 1)
	assert.Equal(t, "/api/*", loadedPluginConfig.ServiceGateway.PathPrefix[0])
}

func TestHandleStrictStaticContentPath(t *testing.T) {
	contentPath := "/hello-world.txt"
	ctx := t.Context()
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	if err != nil {
		t.Error(err)
	}
	testRegisterer := registerer(pluginName)
	staticContentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == contentPath && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = io.WriteString(w, "Hello World!")
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer staticContentServer.Close()

	staticContentServerAddress := staticContentServer.URL
	cfg[pluginName].(map[string]interface{})["static"].([]interface{})[0].(map[string]interface{})["service-host"] = staticContentServerAddress

	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer gatewayServer.Close()
	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	req := httptest.NewRequest(http.MethodGet, contentPath, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello World!", rr.Body.String())
}

func TestHandleGatewayPathPreference(t *testing.T) {
	ctx := t.Context()
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	assert.NoError(t, err)

	testRegisterer := registerer(pluginName)
	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"gateway":"response"}`)
	}))
	defer gatewayServer.Close()

	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "gateway")
}

func TestStaticContentWithMultipleConfigs(t *testing.T) {
	multiStaticConfig := "{" +
		"\"extra_config\": " +
		"{" +
		"\"plugin/http-server\":" +
		"{" +
		"\"name\":[\"hog-static-content\"]," +
		"\"hog-static-content\":" +
		"{" +
		"\"static\":" +
		"[" +
		"{" +
		"\"path-prefix\": \"/assets/*\"," +
		"\"service-host\": \"http://assets-server\"" +
		"}," +
		"{" +
		"\"path-prefix\": \"/images/*\"," +
		"\"service-host\": \"http://images-server\"" +
		"}" +
		"]," +
		"\"service-gateway\": " +
		"{" +
		"\"path-prefix\": [\"/api/*\"]" +
		"}" +
		"}" +
		"}" +
		"}" +
		"}"

	cfg, err := unmarshalExtraConfig([]byte(multiStaticConfig))
	assert.NoError(t, err)

	var loadedPluginConfig PluginConfig
	loadedPluginConfig, err = loadPluginConfig(cfg)

	assert.NoError(t, err)
	assert.Len(t, loadedPluginConfig.Static, 2)
	assert.Equal(t, "/assets/*", loadedPluginConfig.Static[0].PathPrefix)
	assert.Equal(t, "/images/*", loadedPluginConfig.Static[1].PathPrefix)
}

func TestLoadPluginConfigError(t *testing.T) {
	invalidCfg := map[string]interface{}{
		pluginName: "invalid-structure",
	}

	_, err := loadPluginConfig(invalidCfg)
	assert.Error(t, err)
}

func TestRegisterLogger(t *testing.T) {
	testRegisterer := registerer(pluginName)
	mockLogger := &mockLogger{}

	testRegisterer.RegisterLogger(mockLogger)
	assert.NotNil(t, logger)
}

func TestRegisterLoggerWithInvalidType(t *testing.T) {
	testRegisterer := registerer(pluginName)
	originalLogger := logger

	testRegisterer.RegisterLogger("not-a-logger")
	assert.Equal(t, originalLogger, logger)
}

func TestProxyErrorHandling(t *testing.T) {
	ctx := t.Context()
	configWithInvalidHost := "{" +
		"\"extra_config\": " +
		"{" +
		"\"plugin/http-server\":" +
		"{" +
		"\"name\":[\"hog-static-content\"]," +
		"\"hog-static-content\":" +
		"{" +
		"\"static\":" +
		"[" +
		"{" +
		"\"path-prefix\": \"/*\"," +
		"\"service-host\": \"http://non-existent-server-12345.local\"" +
		"}" +
		"]," +
		"\"service-gateway\": " +
		"{" +
		"\"path-prefix\": [\"/api/*\"]" +
		"}" +
		"}" +
		"}" +
		"}" +
		"}"

	cfg, err := unmarshalExtraConfig([]byte(configWithInvalidHost))
	assert.NoError(t, err)

	testRegisterer := registerer(pluginName)
	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer gatewayServer.Close()

	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestProxyWithForwardedHeaders(t *testing.T) {
	ctx := t.Context()
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	assert.NoError(t, err)

	var receivedHeaders http.Header
	staticContentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer staticContentServer.Close()

	cfg[pluginName].(map[string]interface{})["static"].([]interface{})[0].(map[string]interface{})["service-host"] = staticContentServer.URL

	testRegisterer := registerer(pluginName)
	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer gatewayServer.Close()

	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(t, receivedHeaders.Get("X-Forwarded-For"))
	assert.Equal(t, "example.com", receivedHeaders.Get("X-Forwarded-Host"))
	assert.Equal(t, "http", receivedHeaders.Get("X-Forwarded-Proto"))
}

func TestProxyPreservesExistingForwardedHeaders(t *testing.T) {
	ctx := t.Context()
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	assert.NoError(t, err)

	var receivedHeaders http.Header
	staticContentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer staticContentServer.Close()

	cfg[pluginName].(map[string]interface{})["static"].([]interface{})[0].(map[string]interface{})["service-host"] = staticContentServer.URL

	testRegisterer := registerer(pluginName)
	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer gatewayServer.Close()

	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "original.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Contains(t, receivedHeaders.Get("X-Forwarded-For"), "10.0.0.1")
	assert.Contains(t, receivedHeaders.Get("X-Forwarded-For"), "192.168.1.1")
	assert.Equal(t, "https", receivedHeaders.Get("X-Forwarded-Proto"))
	assert.Equal(t, "original.com", receivedHeaders.Get("X-Forwarded-Host"))
}

func TestMatchesWildcard(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wildcard string
		expected bool
	}{
		{"exact match", "/api/users", "/api/*", true},
		{"deep path match", "/api/users/123", "/api/*", true},
		{"root wildcard", "/anything", "/*", true},
		{"no match", "/api/users", "/static/*", false},
		{"no wildcard", "/api/users", "/api/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesWildcard(tt.path, tt.wildcard)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleNoMatchingRoute(t *testing.T) {
	ctx := t.Context()
	cfg, err := unmarshalExtraConfig([]byte(validExtraConfig))
	assert.NoError(t, err)

	// Override static config to use specific paths only
	specificPathConfig := "{" +
		"\"extra_config\": " +
		"{" +
		"\"plugin/http-server\":" +
		"{" +
		"\"name\":[\"hog-static-content\"]," +
		"\"hog-static-content\":" +
		"{" +
		"\"static\":" +
		"[" +
		"{" +
		"\"path-prefix\": \"/static/*\"," +
		"\"service-host\": \"http://static-server\"" +
		"}" +
		"]," +
		"\"service-gateway\": " +
		"{" +
		"\"path-prefix\": [\"/api/*\"]" +
		"}" +
		"}" +
		"}" +
		"}" +
		"}"

	cfg, err = unmarshalExtraConfig([]byte(specificPathConfig))
	assert.NoError(t, err)

	testRegisterer := registerer(pluginName)
	gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer gatewayServer.Close()

	handler, err := testRegisterer.registerHandlers(ctx, cfg, gatewayServer.Config.Handler)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/unmatched-path", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	// Since no route matches, the handler returns without setting a status
	// which defaults to 200 OK from httptest.ResponseRecorder
}

func unmarshalExtraConfig(bytes []byte) (map[string]interface{}, error) {
	var cfg map[string]interface{}
	err := json.Unmarshal(bytes, &cfg)
	if err != nil {
		return cfg, err
	}
	pluginConfigs := make(map[string]interface{})
	extraConfig := cfg["extra_config"].(map[string]interface{})
	httpPluginConfigs := extraConfig["plugin/http-server"].(map[string]interface{})
	names := httpPluginConfigs["name"].([]interface{})
	for _, name := range names {
		key := name.(string)
		pluginConfigs[key] = httpPluginConfigs[key].(map[string]interface{})
	}
	return pluginConfigs, nil
}

// Mock logger for testing
type mockLogger struct {
	debugLogs    []string
	infoLogs     []string
	warningLogs  []string
	errorLogs    []string
	criticalLogs []string
	fatalLogs    []string
}

func (m *mockLogger) Debug(v ...interface{}) { m.debugLogs = append(m.debugLogs, fmt.Sprint(v...)) }
func (m *mockLogger) Info(v ...interface{})  { m.infoLogs = append(m.infoLogs, fmt.Sprint(v...)) }
func (m *mockLogger) Warning(v ...interface{}) {
	m.warningLogs = append(m.warningLogs, fmt.Sprint(v...))
}
func (m *mockLogger) Error(v ...interface{}) { m.errorLogs = append(m.errorLogs, fmt.Sprint(v...)) }
func (m *mockLogger) Critical(v ...interface{}) {
	m.criticalLogs = append(m.criticalLogs, fmt.Sprint(v...))
}
func (m *mockLogger) Fatal(v ...interface{}) { m.fatalLogs = append(m.fatalLogs, fmt.Sprint(v...)) }
