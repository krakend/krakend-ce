package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector"
	"github.com/devopsfaith/krakend-botdetector/krakend"
	"github.com/devopsfaith/krakend-jose"
	muxjose "github.com/devopsfaith/krakend-jose/mux"
	lua "github.com/devopsfaith/krakend-lua/router/mux"
	metrics "github.com/devopsfaith/krakend-metrics/mux"
	opencensus "github.com/devopsfaith/krakend-opencensus/router/mux"
	krakendrate "github.com/devopsfaith/krakend-ratelimit"
	"github.com/devopsfaith/krakend-ratelimit/juju"
	jujurouter "github.com/devopsfaith/krakend-ratelimit/juju/router"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/httptreemux"
	"github.com/devopsfaith/krakend/router/mux"
	router "github.com/devopsfaith/krakend/router/mux"
	"net/http"
	"strings"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	handlerFactory := RateLimitHandlerFactory
	handlerFactory = lua.HandlerFactory(logger, handlerFactory, httptreemux.ParamsExtractor)
	handlerFactory = muxjose.HandlerFactory(handlerFactory, httptreemux.ParamsExtractor, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	handlerFactory = NewBotDetector(handlerFactory, logger)
	return handlerFactory
}

type handlerFactory struct{}

// New checks the configuration and, if required, wraps the handler factory with a bot detector middleware
func NewBotDetector(hf mux.HandlerFactory, l logging.Logger) mux.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) http.HandlerFunc {
		next := hf(cfg, p)

		detectorCfg, err := krakend.ParseConfig(cfg.ExtraConfig)
		if err == krakend.ErrNoConfig {
			l.Debug("botdetector: ", err.Error())
			return next
		}
		if err != nil {
			l.Warning("botdetector: ", err.Error())
			return next
		}

		d, err := botdetector.New(detectorCfg)
		if err != nil {
			l.Warning("botdetector: unable to create the LRU detector:", err.Error())
			return next
		}
		return BotDetectorhandler(d, next)
	}
}

func BotDetectorhandler(f botdetector.DetectorFunc, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if f(r) {
			http.Error(w, "bot rejected", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// HandlerFactory is the out-of-the-box basic ratelimit handler factory using the default krakend endpoint
// handler for the mux router
var RateLimitHandlerFactory = NewRateLimiterMw(mux.CustomEndpointHandler(mux.NewRequestBuilder(httptreemux.ParamsExtractor)))

// NewRateLimiterMw builds a rate limiting wrapper over the received handler factory.
func NewRateLimiterMw(next mux.HandlerFactory) mux.HandlerFactory {
	return func(remote *config.EndpointConfig, p proxy.Proxy) http.HandlerFunc {
		handlerFunc := next(remote, p)

		cfg := jujurouter.ConfigGetter(remote.ExtraConfig).(jujurouter.Config)
		if cfg == jujurouter.ZeroCfg || (cfg.MaxRate <= 0 && cfg.ClientMaxRate <= 0) {
			return handlerFunc
		}

		if cfg.MaxRate > 0 {
			handlerFunc = NewEndpointRateLimiterMw(juju.NewLimiter(float64(cfg.MaxRate), cfg.MaxRate))(handlerFunc)
		}
		if cfg.ClientMaxRate > 0 {
			switch strings.ToLower(cfg.Strategy) {
			case "header":
				handlerFunc = NewHeaderLimiterMw(cfg.Key, float64(cfg.ClientMaxRate), cfg.ClientMaxRate)(handlerFunc)
			}
		}
		return handlerFunc
	}
}

// EndpointMw is a function that decorates the received handlerFunc with some rateliming logic
type EndpointMw func(http.HandlerFunc) http.HandlerFunc

// NewEndpointRateLimiterMw creates a simple ratelimiter for a given handlerFunc
func NewEndpointRateLimiterMw(tb juju.Limiter) EndpointMw {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !tb.Allow() {
				http.Error(w, krakendrate.ErrLimited.Error(), 503)
				return
			}

			next(w, r)
		}
	}
}

// NewHeaderLimiterMw creates a token ratelimiter using the value of a header as a token
func NewHeaderLimiterMw(header string, maxRate float64, capacity int64) EndpointMw {
	return NewTokenLimiterMw(HeaderTokenExtractor(header), juju.NewMemoryStore(maxRate, capacity))
}

// TokenExtractor defines the interface of the functions to use in order to extract a token for each request
type TokenExtractor func(r *http.Request) string

// HeaderTokenExtractor returns a TokenExtractor that looks for the value of the designed header
func HeaderTokenExtractor(header string) TokenExtractor {
	return func(r *http.Request) string { return r.Header.Get(header) }
}

// NewTokenLimiterMw returns a token based ratelimiting endpoint middleware with the received TokenExtractor and LimiterStore
func NewTokenLimiterMw(tokenExtractor TokenExtractor, limiterStore krakendrate.LimiterStore) EndpointMw {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenKey := tokenExtractor(r)
			if tokenKey == "" {
				http.Error(w, krakendrate.ErrLimited.Error(), http.StatusTooManyRequests)
				return
			}
			if !limiterStore(tokenKey).Allow() {
				http.Error(w, krakendrate.ErrLimited.Error(), http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}

func (h handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory) router.HandlerFactory {
	return NewHandlerFactory(l, m, r)
}
