package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/7phs/router/internal/pkg/bridge"
	"github.com/7phs/router/internal/pkg/cache"
	"github.com/7phs/router/internal/pkg/config"
	"github.com/7phs/router/internal/pkg/external_routing_data"
	"github.com/7phs/router/internal/pkg/rest_api"
)

func main() {
	log.Println("router: start")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("failed to load configuration from environmet variable: %v", err)
	}

	dataCache := cache.NewInMemory()
	osmrData := external_routing_data.NewOSMR(cfg.OSRM)
	routingData := bridge.NewBridge(dataCache, osmrData)
	restAPI := rest_api.NewRestAPIServer(cfg.HttpConfig, routingData)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		if err := restAPI.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("failed to start REST API server:", err)
		}
	}()

	select {
	case <-sigs:
	case <-ctx.Done():
	}

	if ctx.Err() == nil {
		if err := restAPI.Shutdown(ctx); err != nil {
			log.Println("router: failed to shutdown REST API:", err)
		}
	}

	log.Println("router: finish")
}
