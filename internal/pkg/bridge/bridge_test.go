package bridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/testutil"
)

type mockCache struct {
	cached pkg.DestinationMeasureList
	stored pkg.DestinationMeasureList
}

func (m *mockCache) GetDestinationMeasures(_ context.Context, _ pkg.Point, _ []pkg.Point) (pkg.DestinationMeasureList, error) {
	return m.cached, nil
}

func (m *mockCache) StoreDestinationMeasures(_ context.Context, _ pkg.Point, destinationMeasure pkg.DestinationMeasureList) error {
	m.stored = destinationMeasure
	return nil
}

type mockExternal struct {
	data  pkg.DestinationMeasureList
	calls int
}

func (m *mockExternal) GetDestinationMeasures(_ context.Context, _ pkg.Point, _ []pkg.Point) (pkg.DestinationMeasureList, error) {
	m.calls++
	return m.data, nil
}

func TestBridge_GetDestinationMeasuresNoCache(t *testing.T) {
	src := testutil.MustParsePoint(t, "13.38886,52.517037")
	dst := testutil.MustParsePoints(t, "13.397634,52.529407", "13.428555,52.523219")
	expectedDestinationList := testutil.GenerateDestinationMeasureList(dst)

	cache := &mockCache{}
	external := &mockExternal{
		data: expectedDestinationList,
	}

	bridge := NewBridge(cache, external)
	actual, err := bridge.GetDestinationMeasures(context.Background(), src, dst)
	require.NoError(t, err)
	require.Equal(t, 1, external.calls)

	assert.Equal(t, expectedDestinationList, actual)
	assert.Equal(t, expectedDestinationList, cache.stored)
}

func TestBridge_GetDestinationMeasuresOnlyCache(t *testing.T) {
	src := testutil.MustParsePoint(t, "13.38886,52.517037")
	dst := testutil.MustParsePoints(t, "13.397634,52.529407", "13.428555,52.523219")
	expectedDestinationList := testutil.GenerateDestinationMeasureList(dst)

	cache := &mockCache{
		cached: expectedDestinationList,
	}
	external := &mockExternal{}

	bridge := NewBridge(cache, external)
	actual, err := bridge.GetDestinationMeasures(context.Background(), src, dst)
	require.NoError(t, err)

	require.Equal(t, 0, external.calls)

	assert.Equal(t, expectedDestinationList, actual)
	assert.Empty(t, cache.stored)
}

func TestBridge_GetDestinationIncludesEquals(t *testing.T) {
	src := testutil.MustParsePoint(t, "13.38886,52.517037")
	dst := testutil.MustParsePoints(t, "13.38886,52.517037", "13.397634,52.529407", "13.428555,52.523219")
	lst := testutil.GenerateDestinationMeasureList(dst)

	cacheDestinationList := lst[1:2]
	extDestinationList := lst[2:]
	expectedDestinationList := append(
		pkg.DestinationMeasureList{pkg.DestinationMeasure{Destination: dst[0]}},
		lst[1:]...,
	)

	cache := &mockCache{
		cached: cacheDestinationList,
	}
	external := &mockExternal{
		data: extDestinationList,
	}

	bridge := NewBridge(cache, external)
	actual, err := bridge.GetDestinationMeasures(context.Background(), src, dst)
	require.NoError(t, err)
	require.Equal(t, 1, external.calls)

	assert.Equal(t, expectedDestinationList, actual)
	assert.Equal(t, extDestinationList, cache.stored)
}
