package routehandler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handle StatsD Gauge request
func (routeHandler *RouteHandler) HandleGaugeRequest(w http.ResponseWriter, r *http.Request) {
	// get key
	vars := mux.Vars(r)
	key := routeHandler.metricPrefix + vars["key"]

	// get gauge shift
	shiftPostFormValue := r.PostFormValue("shift")
	if shiftPostFormValue != "" {
		// get value
		value, err := strconv.Atoi(shiftPostFormValue)
		if err != nil {
			http.Error(w, "Invalid gauge shift specified", 400)
		}
		// send request
		routeHandler.statsdClient.GaugeShift(key, value)
		return
	}

	// get gauge value
	var value = 1
	valuePostFormValue := r.PostFormValue("value")
	if valuePostFormValue != "" {
		// get value
		var err error
		value, err = strconv.Atoi(valuePostFormValue)
		if err != nil {
			http.Error(w, "Invalid gauge value specified", 400)
		}
	}

	// send gauge value request
	routeHandler.statsdClient.Gauge(key, value)

}
