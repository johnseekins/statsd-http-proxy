package router

import (
	"net/http"

	"github.com/GoMetric/statsd-http-proxy/proxy/middleware"
	"github.com/GoMetric/statsd-http-proxy/proxy/routehandler"
	"github.com/gorilla/mux"
)

// NewGorillaMuxRouter creates Gorilla Mux router
func NewGorillaMuxRouter(
	routeHandler *routehandler.RouteHandler,
	tokenSecret string,
) http.Handler {
	// build router
	router := mux.NewRouter().StrictSlash(true)

	// register common middlewares
	router.Use(middleware.ValidateCORS)

	// register http request handlers
	router.Handle(
		"/heartbeat",
		http.HandlerFunc(routeHandler.HandleHeartbeatRequest),
	)

	router.Handle(
		"/count/{key}",
		middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleCountRequest), tokenSecret),
	).Methods(http.MethodPost, http.MethodOptions)

	router.Handle(
		"/gauge/{key}",
		middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleGaugeRequest), tokenSecret),
	).Methods(http.MethodPost, http.MethodOptions)

	router.Handle(
		"/timing/{key}",
		middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleTimingRequest), tokenSecret),
	).Methods(http.MethodPost, http.MethodOptions)

	router.Handle(
		"/set/{key}",
		middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleSetRequest), tokenSecret),
	).Methods(http.MethodPost, http.MethodOptions)

	return router
}
