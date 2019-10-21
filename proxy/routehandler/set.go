package routehandler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handle StatsD Set request
func (routeHandler *routeHandler) handleSetRequest(w http.ResponseWriter, r *http.Request) {
	// get key
	vars := mux.Vars(r)
	key := routeHandler.metricPrefix + vars["key"]

	// get set value
	var value = 1
	valuePostFormValue := r.PostFormValue("value")
	if valuePostFormValue != "" {
		var err error
		value, err = strconv.Atoi(valuePostFormValue)
		if err != nil {
			http.Error(w, "Invalid set value specified", 400)
		}
	}

	// send request
	routeHandler.statsdClient.Set(key, value)
}
