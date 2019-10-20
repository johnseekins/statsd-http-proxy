package routehandler

import (
	"net/http"

	"github.com/GoMetric/go-statsd-client"
	"github.com/gorilla/mux"
	"github.com/GoMetric/statsd-http-proxy/proxy/middleware"
)

type routeHandler struct {
	statsdClient *statsd.Client
	metricPrefix string
}

// NewRouter creates router
func NewRouter(
	statsdClient *statsd.Client,
	metricPrefix string,
	tokenSecret string,
) http.Handler {
	// build router
	router := mux.NewRouter().StrictSlash(true)

	// build route handler
	routeHandler := routeHandler{
		statsdClient,
		metricPrefix,
	}

	// register http request handlers
	router.Handle(
		"/heartbeat",
		middleware.ValidateCORS(http.HandlerFunc(routeHandler.handleHeartbeatRequest))).Methods("GET")

	router.Handle(
		"/count/{key}",
		middleware.ValidateCORS(
			middleware.ValidateJWT(
				http.HandlerFunc(routeHandler.handleCountRequest),
				tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/gauge/{key}",
		middleware.ValidateCORS(
			middleware.ValidateJWT(
				http.HandlerFunc(routeHandler.handleGaugeRequest),
				tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/timing/{key}",
		middleware.ValidateCORS(
			middleware.ValidateJWT(
				http.HandlerFunc(routeHandler.handleTimingRequest),
				tokenSecret)),
	).Methods("POST")

	router.Handle(
		"/set/{key}",
		middleware.ValidateCORS(
			middleware.ValidateJWT(
				http.HandlerFunc(routeHandler.handleSetRequest),
				tokenSecret)),
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
