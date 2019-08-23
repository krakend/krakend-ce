package krakend

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"encoding/json"
	"crypto/sha256"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend-jose"
	"github.com/geetarista/go-bloomd/bloomd"
)

const Namespace = "github.com/openrm/bloomd"

const (
	claimIssuedAt = "iat"
	claimExpirationTime = "exp"
)

var hashFields = []string{"id", "organizationId", claimIssuedAt, claimExpirationTime}

// errors
var (
	errNoConfig = errors.New("no config for bloomd")
	errInvalidConfig = errors.New("invalid config for bloomd")
	errNoFilterName = errors.New("filter name is required")
	errFieldNotExist = errors.New("token missing required field")
	errInvalidField = errors.New("token contains invalid field")
)

// jose.Rejecter implementation
type rejecter struct {
	filter *bloomd.Filter
	logger logging.Logger
}

func (r rejecter) assertFields(claims map[string]interface{}) ([]string, error) {
	fields := make([]string, len(hashFields))
	for i, k := range hashFields {
		v, ok := claims[k]
		if !ok {
			return fields, errFieldNotExist
		}
		str, ok := v.(string)
		if !ok {
			return fields, errInvalidField
		}
		fields[i] = str
	}
	return fields, nil
}

func (r rejecter) calcHash(fields []string) string {
	id := strings.Join(fields, ".")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (r rejecter) Reject(claims map[string]interface{}) bool {
	if r.filter == nil || r.filter.Conn == nil {
		return true
	}
	fields, err := r.assertFields(claims)
	if err != nil {
		return false
	}
	hash := r.calcHash(fields)
	found, err := r.filter.Multi([]string{hash})
	if err != nil {
		r.logger.Error("Bloomd error:", err.Error())
	}
	if len(found) > 0 && found[0] {
		return false
	}
	return true
}

type nopRejecter struct {}
func (nr nopRejecter) Reject(map[string]interface{}) bool { return true }


// config map
type bloomdConfig struct {
	Name string `json:"name"`
	Address string `json:"server_addr"`
}


func createFilter(addr string, filter *bloomd.Filter) error {
	client := bloomd.NewClient(addr)
	return client.CreateFilter(filter)
}

func registerBloomd(scfg config.ServiceConfig, logger logging.Logger) (jose.Rejecter, error) {
	data, ok := scfg.ExtraConfig[Namespace]
	if !ok {
		logger.Debug(errNoConfig.Error())
		return nopRejecter{}, errNoConfig
	}
	raw, err := json.Marshal(data)
	if err != nil {
		logger.Debug(errInvalidConfig.Error())
		return nopRejecter{}, errInvalidConfig
	}
	var cfg bloomdConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		logger.Debug(err.Error(), string(raw))
		return nopRejecter{}, errInvalidConfig
	}
	if cfg.Name == "" {
		return nopRejecter{}, errNoFilterName
	}
	filter := &bloomd.Filter{ Name: cfg.Name }
	if err := createFilter(cfg.Address, filter); err != nil {
		log.Fatalf("Bloomd filter creation failed (%s):", cfg.Address, err.Error())
	}
	logger.Info("BLOOMD: connected")
	return rejecter{filter, logger}, nil
}
