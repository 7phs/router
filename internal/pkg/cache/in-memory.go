package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/bridge"
)

var (
	_ bridge.Cache = (*InMemoryCache)(nil)
)

type InMemoryCacheValue struct {
	measure      pkg.DestinationMeasure
	recentUsedAt time.Time
}

type InMemoryCache struct {
	lock      sync.RWMutex
	cacheData map[string]*InMemoryCacheValue
}

func NewInMemory() *InMemoryCache {
	return &InMemoryCache{
		cacheData: make(map[string]*InMemoryCacheValue),
	}
}

func (c *InMemoryCache) GetDestinationMeasures(_ context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error) {
	if len(dst) == 0 {
		return nil, nil
	}

	var (
		result = pkg.NewDestinationMeasureList(dst)
		now    = time.Now()
	)

	c.lock.RLock()
	defer c.lock.RUnlock()

	for i, dstPoint := range dst {
		key := cacheKey(src, dstPoint)
		value, ok := c.cacheData[key]
		if !ok {
			continue
		}

		result[i] = value.measure
		value.recentUsedAt = now
	}

	return result, nil
}

func (c *InMemoryCache) StoreDestinationMeasures(_ context.Context, src pkg.Point, destinationMeasure pkg.DestinationMeasureList) error {
	if len(destinationMeasure) == 0 {
		return nil
	}

	log.Println("cache: stores", src, " + number of destination measures:", len(destinationMeasure))
	now := time.Now()

	c.lock.Lock()
	defer c.lock.Unlock()

	for _, measure := range destinationMeasure {
		// TODO: needs to check it it makes sense to cache some type of errors as a preventing measure of repeatable requests to an external source
		if measure.Err != nil {
			continue
		}

		key := cacheKey(src, measure.Destination)
		value, ok := c.cacheData[key]
		if !ok {
			value = &InMemoryCacheValue{}
			c.cacheData[key] = value
		}

		value.measure = measure
		value.recentUsedAt = now
	}

	return nil
}

// TODO: cleanup cache based on limitation TTL, number of cached records

func cacheKey(point1, point2 pkg.Point) string {
	// travel from one point and back can require different time of travel
	return point1.Encoded + "|" + point2.Encoded
}
