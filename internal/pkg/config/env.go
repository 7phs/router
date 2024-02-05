package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	DefaultHttpPort       = 80
	DefaultHttpTimeout    = 30 * time.Second
	DefaultOSMRLimit      = 50
	DefaultOSMRHost       = "https://router.project-osrm.org"
	DefaultRequestTimeout = 20 * time.Second
)

func LoadFromEnv() (Config, error) {
	var resultErr error

	httpConfig, err := LoadHttpFromEnv()
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	osrmConfig, err := LoadOSRMFromEnv()
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	if resultErr != nil {
		return Config{}, resultErr
	}

	return Config{
		HttpConfig: httpConfig,
		OSRMConfig: osrmConfig,
	}, nil
}

func LoadHttpFromEnv() (HttpConfig, error) {
	var resultErr error

	httpPort, err := getInt("ROUTER_HTTP_PORT", DefaultHttpPort)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	httpTimeout, err := getDuration("ROUTER_HTTP_TIMEOUT", DefaultHttpTimeout)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	if resultErr != nil {
		return HttpConfig{}, resultErr
	}

	return HttpConfig{
		Port:    httpPort,
		Timeout: httpTimeout,
	}, nil
}

func LoadOSRMFromEnv() (OSRMConfig, error) {
	var resultErr error

	osrmHost, err := getString("ROUTER_OSRM_HOST", DefaultOSMRHost)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	osrmLimit, err := getInt("ROUTER_OSRM_REQUEST_LIMIT", DefaultOSMRLimit)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	osrmRequestTimeout, err := getDuration("ROUTER_OSRM_REQUEST_TIMEOUT", DefaultRequestTimeout)
	if err != nil {
		resultErr = errors.Join(resultErr, err)
	}

	if resultErr != nil {
		return OSRMConfig{}, resultErr
	}

	return OSRMConfig{
		Host:                 osrmHost,
		LimitRequestsPerTime: osrmLimit,
		RequestTimeout:       osrmRequestTimeout,
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
