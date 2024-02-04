package bridge

import (
	"context"
	"fmt"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/rest_api"
)

var (
	_ rest_api.RoutingData = (*Bridge)(nil)
)

type Cache interface {
	GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, error)
	StoreDestinationMeasures(ctx context.Context, src pkg.Point, destinationMeasure []pkg.DestinationMeasure) error
}

type ExternalRoutingData interface {
	GetDestinationMeasure(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, error)
}

type Bridge struct {
	cache               Cache
	externalRoutingData ExternalRoutingData
}

func NewBridge(cache Cache, externalRoutingData ExternalRoutingData) *Bridge {
	return &Bridge{
		cache:               cache,
		externalRoutingData: externalRoutingData,
	}
}

func (b *Bridge) GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, error) {
	measures, notCachedDst, err := b.fetchDataFromCache(ctx, src, dst)
	if err != nil {
		return nil, err
	}

	if len(notCachedDst) == 0 {
		return measures, nil
	}

	newMeasures, err := b.fetchDataFromExternalRoutingData(ctx, src, notCachedDst)
	if err != nil {
		return measures, err
	}

	err = b.cache.StoreDestinationMeasures(ctx, src, newMeasures)
	if err != nil {
		return measures, err
	}

	for _, newData := range newMeasures {
		if newData.Destination.Index >= 0 && newData.Destination.Index < len(measures) {
			measures[newData.Destination.Index] = newData
		}
	}

	return measures, nil
}

func (b *Bridge) fetchDataFromCache(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, []pkg.Point, error) {
	measures := pkg.NewDestinationMeasureList(len(dst))
	cachedData, err := b.cache.GetDestinationMeasures(ctx, src, dst)
	if err != nil {
		return nil, nil, err
	}
	if len(cachedData) != len(measures) {
		return nil, nil, fmt.Errorf("unexpected data from cache, len of result not equal number of destination points")
	}

	var notCachedDst []pkg.Point

	for i, data := range cachedData {
		if data.Err != nil {
			dst[i].Index = i
			notCachedDst = append(notCachedDst, dst[i])
			continue
		}

		measures[i] = data
	}

	return measures, notCachedDst, nil
}

func (b *Bridge) fetchDataFromExternalRoutingData(ctx context.Context, src pkg.Point, dst []pkg.Point) ([]pkg.DestinationMeasure, error) {
	return b.externalRoutingData.GetDestinationMeasure(ctx, src, dst)
}
