package routehandler

import (
	"fmt"
	"net/http"

	"github.com/johnseekins/statsd-http-proxy/proxy/statsdclient"
)

// RouteHandler as a collection of route handlers
type RouteHandler struct {
	statter statsdclient.Statter
}

// NewRouteHandler creates collection of route handlers
func NewRouteHandler(
	statter statsdclient.Statter,
) *RouteHandler {
	// build route handler
	routeHandler := RouteHandler{
		statter,
	}

	return &routeHandler
}

func (routeHandler *RouteHandler) HandleMetric(
	w http.ResponseWriter,
	r *http.Request,
	metricType string,
	metricKey string,
) {
	switch metricType {
	case "count":
		routeHandler.handleCountRequest(w, r, metricKey)
	case "incr":
		routeHandler.handleIncrementRequest(w, r, metricKey)
	case "gauge":
		routeHandler.handleGaugeRequest(w, r, metricKey)
	case "timing":
		routeHandler.handleTimingRequest(w, r, metricKey)
	case "uniq":
		routeHandler.handleUniqueRequest(w, r, metricKey)
	}
}

func (routeHandler *RouteHandler) HandleHeartbeatRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
