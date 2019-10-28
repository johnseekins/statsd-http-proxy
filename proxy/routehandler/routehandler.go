package routehandler

import (
	"net/http"

	"github.com/GoMetric/statsd-http-proxy/proxy/statsdclient"
)

// RouteHandler as a collection of route handlers
type RouteHandler struct {
	statsdClient statsdclient.StatsdClientInterface
	metricPrefix string
}

// NewRouteHandler creates collection of route handlers
func NewRouteHandler(
	statsdClient statsdclient.StatsdClientInterface,
	metricPrefix string,
) *RouteHandler {
	// build route handler
	routeHandler := RouteHandler{
		statsdClient,
		metricPrefix,
	}

	return &routeHandler
}

// GetFullyQualifiedMetricKey return metric key with passed suffix and pre-configured prefix
func (routeHandler *RouteHandler) getFullyQualifiedMetricKey(metricKeySuffix string) string {
	return routeHandler.metricPrefix + metricKeySuffix
}

//HandleMetric reads count, gauge, timing and set metrics from HTTP and sent them to StatsD
func (routeHandler *RouteHandler) HandleMetric(
	w http.ResponseWriter,
	r *http.Request,
	metricType string,
	metricKeySuffix string,
) {
	// get fully qualified metric key
	metricKey := routeHandler.getFullyQualifiedMetricKey(metricKeySuffix)

	// run handler
	switch metricType {
	case "count":
		routeHandler.handleCountRequest(w, r, metricKey)
	case "gauge":
		routeHandler.handleGaugeRequest(w, r, metricKey)
	case "timing":
		routeHandler.handleTimingRequest(w, r, metricKey)
	case "set":
		routeHandler.handleSetRequest(w, r, metricKey)
	}
}
