package krakend

import (
	"fmt"
	"errors"
	"strings"
	"encoding/json"
	"crypto/sha256"
	"time"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend-jose"
	"github.com/geetarista/go-bloomd/bloomd"
)

const Namespace = "github_com/openrm/bloomd"

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
			// return fields, errFieldNotExist
		}

		fields[i] = fmt.Sprintf("%v", v)
	}

	return fields, nil
}

func (r rejecter) calcHash(fields []string) string {
	id := strings.Join(fields, ".")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func (r rejecter) Reject(claims map[string]interface{}) bool {

	if r.filter == nil || r.filter.Conn == nil {
		return false
	}

	fields, err := r.assertFields(claims)

	if err != nil {
		return true
	}

	hash := r.calcHash(fields)

	if r.filter.Conn.Socket != nil {
		r.filter.Conn.Socket.SetReadDeadline(time.Now().Add(1 * time.Second))
		r.filter.Conn.Socket.SetWriteDeadline(time.Now().Add(1 * time.Second))
	}

	found, err := r.filter.Multi([]string{hash})

	if err != nil {
		r.logger.Error("Bloomd error:", err.Error())
		connectAndConfigure(r.filter)
	}

	if len(found) > 0 && found[0] {
		return true
	}

	return false
}

type nopRejecter struct {}
func (nr nopRejecter) Reject(map[string]interface{}) bool { return false }


// config map
type bloomdConfig struct {
	Name string `json:"name"`
	Address string `json:"server_addr"`
}

func connectAndConfigure(filter *bloomd.Filter) bool {
	if filter.Conn.Socket != nil {
		filter.Conn.Socket.Close()
		filter.Conn.Socket = nil
	}

	info, err := filter.Info()

	if err != nil {
		fmt.Println("Error connecting to bloomd:", err)
		return false
	}

	fmt.Println("Connected to bloomd: %v", info)

	filter.Conn.Socket.SetKeepAlive(true)
	filter.Conn.Socket.SetKeepAlivePeriod(20 * time.Second)

	return true
}


func createFilter(addr string, filterName string) *bloomd.Filter {
	client := bloomd.NewClient(addr)
	filter := client.GetFilter(filterName)

	connectAndConfigure(filter)

	return filter
}

func RegisterBloomd(scfg config.ServiceConfig, logger logging.Logger) (jose.Rejecter, error) {

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

	filter := createFilter(cfg.Address, cfg.Name )

	return rejecter{filter, logger}, nil
}
