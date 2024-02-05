package external_routing_data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/config"
	"github.com/7phs/router/internal/pkg/testutil"
)

func TestOSMR_GetDestinationMeasure(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadOSRMFromEnv()
	require.NoError(t, err)

	osrmData := NewOSMR(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	src := testutil.MustParsePoint(t, "13.38886,52.517037")
	dst := testutil.MustParsePoints(t, "13.397634,52.529407", "13.428555,52.523219")

	actualMeasures, err := osrmData.GetDestinationMeasures(ctx, src, dst)
	require.NoError(t, err)

	expectedMeasures := pkg.DestinationMeasureList{
		{
			Destination: dst[0],
			Duration:    260.3,
			Distance:    1886.3,
			Err:         nil,
		},
		{
			Destination: dst[1],
			Duration:    389.3,
			Distance:    3804.2,
			Err:         nil,
		},
	}
	assert.Equal(t, expectedMeasures, actualMeasures)
}
