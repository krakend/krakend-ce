package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector"
	"github.com/devopsfaith/krakend-botdetector/krakend"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/mux"
	"net/http"
)

// Register checks the configuration and, if required, registers a bot detector middleware at the gin engine
func Register(cfg config.ExtraConfig, l logging.Logger, mw []mux.HandlerMiddleware) []mux.HandlerMiddleware {
	detectorCfg, err := krakend.ParseConfig(cfg)
	if err == krakend.ErrNoConfig {
		l.Debug("botdetector middleware: ", err.Error())
		return mw
	}
	if err != nil {
		l.Warning("botdetector middleware: ", err.Error())
		return mw
	}
	d, err := botdetector.New(detectorCfg)
	if err != nil {
		l.Warning("botdetector middleware: unable to createt the LRU detector:", err.Error())
		return mw
	}
	return append(mw, middleware(d))
}

type Middleware struct {
	f botdetector.DetectorFunc
}

func (hm *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hm.f(r) {
			http.Error(w, "bot rejected", http.StatusForbidden)
		}

		h.ServeHTTP(w, r)
	})
}

func middleware(f botdetector.DetectorFunc) mux.HandlerMiddleware {
	return &Middleware{f: f}
}

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
