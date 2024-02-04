package config

import (
	"fmt"
	"time"
)

type HttpConfig struct {
	Port    int
	Timeout time.Duration
}

func (h HttpConfig) Address() string {
	return fmt.Sprintf(":%d", h.Port)
}

type OSRM struct {
	Host                 string
	LimitRequestsPerTime int
}

type Config struct {
	HttpConfig
	OSRM
}