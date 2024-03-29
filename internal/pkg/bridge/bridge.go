package bridge

import (
	"context"
	"log"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/rest_api"
)

var (
	_ rest_api.RoutingData = (*Bridge)(nil)
)

type Cache interface {
	GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error)
	StoreDestinationMeasures(ctx context.Context, src pkg.Point, destinationMeasure pkg.DestinationMeasureList) error
}

type ExternalRoutingData interface {
	GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error)
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

func (b *Bridge) GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error) {
	measures := pkg.NewDestinationMeasureList(dst)

	measures.LabelEqualDestinations(src)

	dst = measures.NotProcessedPoints()
	if len(dst) == 0 {
		return measures, nil
	}

	cachedMeasures, err := b.cache.GetDestinationMeasures(ctx, src, dst)
	if err != nil {
		log.Println("failed to fetch measures from cache:", err)
	}

	if len(cachedMeasures) > 0 {
		measures.UpdateMeasures(cachedMeasures)
	}

	dst = measures.NotProcessedPoints()
	if len(dst) == 0 {
		return measures, nil
	}

	newMeasures, err := b.externalRoutingData.GetDestinationMeasures(ctx, src, dst)
	if err != nil {
		log.Println("failed to fetch measures from external sources:", err)
	}

	if len(newMeasures) > 0 {
		err = b.cache.StoreDestinationMeasures(ctx, src, newMeasures)
		if err != nil {
			log.Println("failed to store new measures into cache:", err)
		}
	}

	measures.UpdateMeasures(newMeasures)

	return measures, nil
}
