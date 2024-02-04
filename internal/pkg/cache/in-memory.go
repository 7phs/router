package cache

import (
	"context"
	"os"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/bridge"
)

var (
	_ bridge.Cache = (*InMemoryCache)(nil)
)

type InMemoryCache struct {
}

func NewInMemory() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, error) {
	return nil, os.ErrInvalid
}

func (c *InMemoryCache) StoreDestinationMeasures(ctx context.Context, src pkg.Point, destinationMeasure []pkg.DestinationMeasure) error {
	return os.ErrInvalid
}
