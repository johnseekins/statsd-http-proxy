package routehandler

import (
	"github.com/GoMetric/go-statsd-client"
)

// RouteHandler as a collection of route handlers
type RouteHandler struct {
	statsdClient *statsd.Client
	metricPrefix string
}

// NewRouteHandler creates collection of route handlers
func NewRouteHandler(
	statsdClient *statsd.Client,
	metricPrefix string,
) *RouteHandler {
	// build route handler
	routeHandler := RouteHandler{
		statsdClient,
		metricPrefix,
	}

	return &routeHandler
}
