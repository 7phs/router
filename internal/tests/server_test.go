package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/7phs/router/internal/pkg/bridge"
	"github.com/7phs/router/internal/pkg/cache"
	"github.com/7phs/router/internal/pkg/config"
	"github.com/7phs/router/internal/pkg/external_routing_data"
	"github.com/7phs/router/internal/pkg/rest_api"
)

func TestRestAPIServer_Routes(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadFromEnv()
	require.NoError(t, err)

	brd := bridge.NewBridge(cache.NewInMemory(), external_routing_data.NewOSMR(cfg.OSRMConfig))
	srv := rest_api.NewRestAPIServer(cfg.HttpConfig, brd)

	server := httptest.NewServer(srv.Handler())
	defer server.Close()

	requestURL, err := url.Parse(server.URL + "/routes?src=13.388860,52.517037&dst=13.428555,52.523219&dst=13.397634,52.529407")
	require.NoError(t, err)

	client := server.Client()
	res, err := client.Get(requestURL.String())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var actual rest_api.RoutesResponse

	err = decoder.Decode(&actual)
	require.NoError(t, err)

	expected := rest_api.RoutesResponse{
		Source: "13.388860,52.517037",
		Routes: []rest_api.Route{
			{
				Destination: "13.397634,52.529407",
				Duration:    260.3,
				Distance:    1886.3,
			},
			{
				Destination: "13.428555,52.523219",
				Duration:    389.3,
				Distance:    3804.2,
			},
		},
	}

	assert.Equal(t, expected, actual)
}
