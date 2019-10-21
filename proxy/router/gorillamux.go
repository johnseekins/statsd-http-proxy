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

	// register http request handlers
	router.Handle(
		"/heartbeat",
		middleware.ValidateCORS(http.HandlerFunc(routeHandler.HandleHeartbeatRequest)),
	).Methods("GET")

	router.Handle(
		"/count/{key}",
		middleware.ValidateCORS(middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleCountRequest), tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/gauge/{key}",
		middleware.ValidateCORS(middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleGaugeRequest), tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/timing/{key}",
		middleware.ValidateCORS(middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleTimingRequest), tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/set/{key}",
		middleware.ValidateCORS(middleware.ValidateJWT(http.HandlerFunc(routeHandler.HandleSetRequest), tokenSecret)),
	).Methods("POST")

	// Handle pre-flight CORS requests
	router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") == "" {
			return
		}

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Headers", middleware.JwtHeaderName+", X-Requested-With, Origin, Accept, Content-Type, Authentication")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
		}
	})

	return router
}
