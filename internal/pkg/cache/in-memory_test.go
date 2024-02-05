package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/testutil"
)

func TestInMemory(t *testing.T) {
	t.Parallel()

	cache := NewInMemory()
	ctx := context.Background()

	dst := testutil.MustParsePoints(t, "13.397634,52.529407", "15.234223,56.219932")
	emptyMeasures := pkg.NewDestinationMeasureList(dst)

	src1 := testutil.MustParsePoint(t, "13.38886,52.517037")
	measuresList1 := testutil.GenerateDestinationMeasureList(dst)

	src2 := testutil.MustParsePoint(t, "13.428555,52.523219")
	measuresList2 := testutil.GenerateDestinationMeasureList(dst)

	// SRC1
	actual, err := cache.GetDestinationMeasures(ctx, src1, dst)
	require.NoError(t, err)
	assert.Equal(t, emptyMeasures, actual)

	require.NoError(t, cache.StoreDestinationMeasures(ctx, src1, measuresList1))

	actual, err = cache.GetDestinationMeasures(ctx, src1, dst)
	require.NoError(t, err)
	assert.Equal(t, measuresList1, actual)

	// SRC2
	actual, err = cache.GetDestinationMeasures(ctx, src2, dst)
	require.NoError(t, err)
	assert.Equal(t, emptyMeasures, actual)

	require.NoError(t, cache.StoreDestinationMeasures(ctx, src2, measuresList2))

	actual, err = cache.GetDestinationMeasures(ctx, src2, dst)
	require.NoError(t, err)
	assert.Equal(t, measuresList2, actual)

	// SRC 1 should be there
	actual, err = cache.GetDestinationMeasures(ctx, src1, dst)
	require.NoError(t, err)
	assert.Equal(t, measuresList1, actual)
}

func TestInMemoryUpdate(t *testing.T) {
	t.Parallel()

	cache := NewInMemory()
	ctx := context.Background()

	src := testutil.MustParsePoint(t, "13.38886,52.517037")
	dst := testutil.MustParsePoints(t, "13.397634,52.529407", "15.234223,56.219932")
	measuresList := testutil.GenerateDestinationMeasureList(dst)
	updatedMeasuresList := testutil.GenerateDestinationMeasureList(dst)[1:]
	expectedMeasureList := append(measuresList[:1], updatedMeasuresList...)

	require.NoError(t, cache.StoreDestinationMeasures(ctx, src, measuresList))

	actual, err := cache.GetDestinationMeasures(ctx, src, dst)
	require.NoError(t, err)
	assert.Equal(t, measuresList, actual)

	require.NoError(t, cache.StoreDestinationMeasures(ctx, src, updatedMeasuresList))

	actual, err = cache.GetDestinationMeasures(ctx, src, dst)
	require.NoError(t, err)
	assert.Equal(t, expectedMeasureList, actual)
}
