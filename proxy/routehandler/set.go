package routehandler

import (
	"net/http"
	"strconv"
)

// Handle StatsD Set request
func (routeHandler *RouteHandler) handleSetRequest(w http.ResponseWriter, r *http.Request, key string) {
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
