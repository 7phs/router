package external_routing_data

import (
	"context"
	"fmt"
	"log"

	osrm "github.com/gojuno/go.osrm"
	geo "github.com/paulmach/go.geo"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/bridge"
	"github.com/7phs/router/internal/pkg/config"
)

var (
	_ bridge.ExternalRoutingData = (*OSMR)(nil)
)

type OSMR struct {
	rateLimit  *rate.Limiter
	osrmClient *osrm.OSRM
}

func NewOSMR(cfg config.OSRMConfig) *OSMR {
	return &OSMR{
		rateLimit:  rate.NewLimiter(rate.Limit(cfg.LimitRequestsPerTime), cfg.LimitRequestsPerTime),
		osrmClient: osrm.NewFromURLWithTimeout(cfg.Host, cfg.RequestTimeout),
	}
}

func (o *OSMR) GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error) {
	var (
		result          = pkg.NewDestinationMeasureList(dst)
		group, groupCtx = errgroup.WithContext(ctx)
	)

	log.Println("external: requests osrm service for", src, dst)

	for i, dstPoint := range dst {
		index := i
		geometry := osrm.NewGeometryFromPath(*geo.NewPath().SetPoints([]geo.Point{src.Point, dstPoint.Point}))

		group.Go(func() error {
			newDst, err := o.fetchRoute(groupCtx, geometry)
			if err == nil {
				newDst.Destination = dst[index]
			}

			result[index] = newDst
			return err
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (o *OSMR) fetchRoute(ctx context.Context, geometry osrm.Geometry) (pkg.DestinationMeasure, error) {
	err := o.rateLimit.Wait(ctx)
	if err != nil {
		return pkg.NewDestinationMeasureErr(err)
	}

	resp, err := o.osrmClient.Route(ctx, osrm.RouteRequest{
		Profile:     "driving",
		Coordinates: geometry,
		Overview:    osrm.OverviewFalse,
	})
	if err != nil {
		return pkg.NewDestinationMeasureErr(err)
	}

	if len(resp.Routes) == 0 {
		return pkg.NewDestinationMeasureErr(fmt.Errorf("routes empty for %v", geometry.String()))
	}

	return pkg.DestinationMeasure{
		Duration: resp.Routes[0].Duration,
		Distance: resp.Routes[0].Distance,
	}, nil
}
