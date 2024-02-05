package testutil

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/7phs/router/internal/pkg"
)

func MustParsePoint(t *testing.T, v string) pkg.Point {
	point, err := pkg.ParsePoint(v, pkg.DefaultFactor)
	require.NoError(t, err)

	return point
}

func MustParsePoints(t *testing.T, v ...string) []pkg.Point {
	var lst []pkg.Point

	for _, s := range v {
		lst = append(lst, MustParsePoint(t, s))
	}

	return lst
}

func GenerateDestinationMeasureList(lst []pkg.Point) pkg.DestinationMeasureList {
	var result pkg.DestinationMeasureList

	for i, point := range lst {
		point.Index = i
		result = append(result, pkg.DestinationMeasure{
			Destination: point,
			Duration:    rand.Float32(),
			Distance:    rand.Float32(),
			Err:         nil,
		})
	}

	return result
}
