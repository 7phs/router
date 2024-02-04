package pkg

import (
	"fmt"
	"strings"

	geo "github.com/paulmach/go.geo"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Point struct {
	Encoded string
	Point   geo.Point
	Index   int
}

func ParsePoint(encoded string) (Point, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return Point{}, fmt.Errorf("failed to parse a coordinare '%s': empty")
	}

	path := geo.NewPathFromEncoding(encoded)

	if len(path.PointSet) != 1 {
		return Point{}, fmt.Errorf("failed to parse a coordinare '%s': number of components not equal 2")
	}

	return Point{
		Encoded: encoded,
		Point:   path.PointSet[0],
	}, nil
}

type DestinationMeasure struct {
	Destination Point
	Duration    float32
	Distance    float32
	Err         error
}

func NewDestinationMeasureErr(err error) (DestinationMeasure, error) {
	return DestinationMeasure{Err: err}, err
}

func NewDestinationMeasureList(sz int) []DestinationMeasure {
	result := make([]DestinationMeasure, sz)

	for i := 0; i < sz; i++ {
		result[i].Err = ErrNotFound
	}

	return result
}
