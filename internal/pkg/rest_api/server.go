package rest_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/gorilla/mux"

	"github.com/7phs/router/internal/pkg"
	"github.com/7phs/router/internal/pkg/config"
)

const MaxHeaderBytes = 1 << 20

type RoutingData interface {
	GetDestinationMeasures(ctx context.Context, src pkg.Point, dst []pkg.Point) (pkg.DestinationMeasureList, error)
}

type RestAPIServer struct {
	routingData RoutingData
	httpServer  *http.Server
}

func NewRestAPIServer(cfg config.HttpConfig, routingData RoutingData) *RestAPIServer {
	api := &RestAPIServer{
		routingData: routingData,
		httpServer: &http.Server{
			Addr:           cfg.Address(),
			ReadTimeout:    cfg.Timeout,
			WriteTimeout:   cfg.Timeout,
			MaxHeaderBytes: MaxHeaderBytes,
		},
	}

	r := mux.NewRouter()
	r.HandleFunc("/routes", api.Routes).
		Methods("GET")

	api.httpServer.Handler = r

	return api
}

func (api *RestAPIServer) Handler() http.Handler {
	return api.httpServer.Handler
}

func (api *RestAPIServer) Start() error {
	return api.httpServer.ListenAndServe()
}

func (api *RestAPIServer) Shutdown(ctx context.Context) error {
	return api.httpServer.Shutdown(ctx)
}

type ErrorResponse struct {
	Err string `json:"error_message"`
}

func parseQueryParams(query url.Values) (pkg.Point, []pkg.Point, error) {
	var (
		resultErr error
		err       error
		src       pkg.Point
		dst       []pkg.Point
	)

	values, ok := query["src"]
	switch {
	case !ok:
		resultErr = errors.Join(resultErr, fmt.Errorf("src: not presented"))
	case len(values) != 1:
		resultErr = errors.Join(resultErr, fmt.Errorf("src: unexpected numbers of values - %v", len(values)))
	default:
		if src, err = pkg.ParsePoint(values[0], pkg.DefaultFactor); err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	}

	values, ok = query["dst"]
	if !ok {
		resultErr = errors.Join(resultErr, fmt.Errorf("dst: not presented"))
	} else {
		for _, value := range values {
			if dstValue, err := pkg.ParsePoint(value, pkg.DefaultFactor); err != nil {
				resultErr = errors.Join(resultErr, err)
			} else {
				dst = append(dst, dstValue)
			}
		}
	}

	return src, dst, resultErr
}

func writeError(status int, w http.ResponseWriter, err error) {
	w.WriteHeader(status)
	b, _ := json.Marshal(ErrorResponse{Err: err.Error()})
	_, _ = w.Write(b)
}

type Route struct {
	Destination string  `json:"destination"`
	Duration    float32 `json:"duration,omitempty"`
	Distance    float32 `json:"distance,omitempty"`
	Err         string  `json:"err,omitempty"`
}

type RoutesResponse struct {
	Source string  `json:"sources"`
	Routes []Route `json:"routes"`
}

func RoutesResponseFromDestinationMeasure(src pkg.Point, measures []pkg.DestinationMeasure) RoutesResponse {
	routes := make([]Route, 0, len(measures))

	for _, measure := range measures {
		var errMsg string

		if measure.Err != nil {
			errMsg = measure.Err.Error()
		}

		routes = append(routes, Route{
			Destination: measure.Destination.Encoded,
			Duration:    measure.Duration,
			Distance:    measure.Distance,
			Err:         errMsg,
		})
	}

	return RoutesResponse{
		Source: src.Encoded,
		Routes: routes,
	}
}

func (api *RestAPIServer) Routes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	src, dst, err := parseQueryParams(r.URL.Query())
	if err != nil {
		writeError(http.StatusBadRequest, w, err)
		return
	}

	destinationMeasures, err := api.routingData.GetDestinationMeasures(r.Context(), src, dst)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err)
		return
	}

	sort.SliceStable(destinationMeasures, func(i, j int) bool {
		if destinationMeasures[i].Duration < destinationMeasures[j].Duration {
			return true
		}

		return destinationMeasures[i].Distance < destinationMeasures[j].Distance
	})

	response := RoutesResponseFromDestinationMeasure(src, destinationMeasures)
	b, err := json.Marshal(response)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err)
		return
	}

	_, _ = w.Write(b)
}
