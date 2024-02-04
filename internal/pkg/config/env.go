package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	DefaultHttpPort    = 80
	DefaultHttpTimeout = 10 * time.Second
	DefaultOSMRLimit   = 50
	DefaultOSMRHost    = "https://router.project-osrm.org"
)

func LoadFromEnv() (Config, error) {
	var resultErr error

	httpPort, err := getInt("ROUTER_HTTP_PORT", DefaultHttpPort)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	httpTimeout, err := getDuration("ROUTER_HTTP_TIMEOUT", DefaultHttpTimeout)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	osrmLimit, err := getInt("ROUTER_OSRM_REQUEST_LIMIT", DefaultOSMRLimit)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	osrmHost, err := getString("ROUTER_OSRM_REQUEST_LIMIT", DefaultOSMRHost)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	if resultErr != nil {
		return Config{}, resultErr
	}

	return Config{
		HttpConfig: HttpConfig{
			Port:    httpPort,
			Timeout: httpTimeout,
		},
		OSRM: OSRM{
			Host:                 osrmHost,
			LimitRequestsPerTime: osrmLimit,
		},
	}, nil
}

func getInt(name string, def int) (int, error) {
	return getValue(name, def, func(s string) (int, error) {
		v, err := strconv.ParseInt(s, 10, 32)
		return int(v), err
	})
}

func getDuration(name string, def time.Duration) (time.Duration, error) {
	return getValue(name, def, time.ParseDuration)
}

func getString(name, def string) (string, error) {
	return getValue(name, def, func(s string) (string, error) {
		return s, nil
	})
}

func getValue[T any](name string, def T, parse func(string) (T, error)) (T, error) {
	vs := os.Getenv(name)
	if vs == "" {
		return def, nil
	}

	v, err := parse(vs)
	if err != nil {
		var noValue T
		return noValue, err
	}

	return v, err
}
