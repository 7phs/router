package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	geo "github.com/paulmach/go.geo"
)

type Factor struct {
	factor         int
	multiplication int64
}

var (
	DefaultFactor = Factor{
		factor:         6,
		multiplication: 1000000,
	}
	ErrNotFound = fmt.Errorf("not found")
)

type Point struct {
	Encoded string
	Point   geo.Point
	Index   int
}

func ParsePoint(encoded string, factor Factor) (Point, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return Point{}, fmt.Errorf("failed to parse a coordinare '%s': empty", encoded)
	}

	parts := strings.SplitN(encoded, ",", 2)
	if len(parts) < 2 {
		return Point{}, fmt.Errorf("failed to parse a coordinare '%s': one part only", encoded)
	}

	component := [2]float64{}

	for i := 0; i < 2; i++ {
		numParts := strings.SplitN(parts[i], ".", 2)
		if len(numParts) < 2 {
			return Point{}, fmt.Errorf("failed to parse a coordinare '%s': invalid coord", encoded)
		}
		i1, err := strconv.ParseInt(numParts[0], 10, 32)
		if err != nil {
			return Point{}, fmt.Errorf("failed to parse a coordinare '%s': invalid coord", encoded)
		}

		if len(numParts[1]) > factor.factor {
			return Point{}, fmt.Errorf("failed to parse a coordinare '%s': a coord has a factor is great than %d", encoded, factor.factor)
		}
		if len(numParts[1]) < factor.factor {
			numParts[1] += strings.Repeat("0", factor.factor-len(numParts[1]))
		}
		i2, err := strconv.ParseInt(numParts[1], 10, 32)
		if err != nil {
			return Point{}, fmt.Errorf("failed to parse a coordinare '%s': invalid coord", encoded)
		}

		component[i] = float64(i1*factor.multiplication+i2) / float64(factor.multiplication)
	}

	return Point{
		Encoded: encoded,
		Point:   geo.NewPathFromXYData([][2]float64{component}).PointSet[0],
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

type DestinationMeasureList []DestinationMeasure

func NewDestinationMeasureList(points []Point) DestinationMeasureList {
	result := make([]DestinationMeasure, len(points))

	for i, point := range points {
		point.Index = i
		result[i] = DestinationMeasure{
			Destination: point,
			Err:         ErrNotFound,
		}
	}

	return result
}

func (l *DestinationMeasureList) LabelEqualDestinations(src Point) {
	for i, measure := range *l {
		if src.Encoded != measure.Destination.Encoded {
			continue
		}

		(*l)[i] = DestinationMeasure{
			Destination: measure.Destination,
			Duration:    0.0,
			Distance:    0.0,
			Err:         nil,
		}
	}
}

func (l DestinationMeasureList) NotProcessedPoints() []Point {
	var notProcessed []Point

	for _, measure := range l {
		if errors.Is(measure.Err, ErrNotFound) {
			notProcessed = append(notProcessed, measure.Destination)
		}
	}

	return notProcessed
}

func (l *DestinationMeasureList) UpdateMeasures(updatedMeasures DestinationMeasureList) {
	for _, updatedData := range updatedMeasures {
		if updatedData.Err != nil {
			continue
		}
		if updatedData.Destination.Index < 0 || updatedData.Destination.Index >= len(*l) {
			continue
		}

		(*l)[updatedData.Destination.Index] = updatedData
	}
}
