package routehandler

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (routeHandler *RouteHandler) handleCountRequest(w http.ResponseWriter, r *http.Request, key string) {
	var value = 0
	valuePostFormValue := r.PostFormValue("n")
	if valuePostFormValue != "" {
		var err error
		value, err = strconv.Atoi(valuePostFormValue)
		if err != nil {
			http.Error(w, "Invalid n specified", 400)
			return
		}
	}

	valuePostFormTags := r.PostFormValue("tags")
	valuePostFormTags = strings.TrimSpace(valuePostFormTags)
	if valuePostFormTags != "" {
		if validateTags(valuePostFormTags) {
			key += "," + valuePostFormTags
		} else {
			http.Error(w, "Invalid tags specified", 400)
			return
		}
	}

	routeHandler.statter.Count(key, value)
}

func (routeHandler *RouteHandler) handleIncrementRequest(w http.ResponseWriter, r *http.Request, key string) {
	valuePostFormValue := r.PostFormValue("n")
	if valuePostFormValue != "" {
		http.Error(w, "Cannot specify n or value for incr", 400)
		return
	}

	valuePostFormTags := r.PostFormValue("tags")
	valuePostFormTags = strings.TrimSpace(valuePostFormTags)
	if valuePostFormTags != "" {
		if validateTags(valuePostFormTags) {
			key += "," + valuePostFormTags
		} else {
			http.Error(w, "Invalid tags specified", 400)
			return
		}
	}

	routeHandler.statter.Increment(key)
}

func (routeHandler *RouteHandler) handleGaugeRequest(w http.ResponseWriter, r *http.Request, key string) {
	var value = 0
	valuePostFormValue := r.PostFormValue("value")
	if valuePostFormValue != "" {
		var err error
		value, err = strconv.Atoi(valuePostFormValue)
		if err != nil {
			http.Error(w, "Invalid gauge value specified", 400)
			return
		}
	}

	valuePostFormTags := r.PostFormValue("tags")
	valuePostFormTags = strings.TrimSpace(valuePostFormTags)
	if valuePostFormTags != "" {
		if validateTags(valuePostFormTags) {
			key += "," + valuePostFormTags
		} else {
			http.Error(w, "Invalid tags specified", 400)
			return
		}
	}

	routeHandler.statter.Gauge(key, value)
}

func (routeHandler *RouteHandler) handleTimingRequest(w http.ResponseWriter, r *http.Request, key string) {
	var duration time.Duration

	valuePostFormDur := r.PostFormValue("dur")
	if valuePostFormDur != "" {
		d, err := strconv.ParseInt(r.PostFormValue("dur"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid dur specified", 400)
			return
		}

		duration = time.Duration(d) * time.Millisecond
	}

	valuePostFormTags := r.PostFormValue("tags")
	valuePostFormTags = strings.TrimSpace(valuePostFormTags)
	if valuePostFormTags != "" {
		if validateTags(valuePostFormTags) {
			key += "," + valuePostFormTags
		} else {
			http.Error(w, "Invalid tags specified", 400)
			return
		}
	}

	// send request
	routeHandler.statter.Timing(key, duration)
}

func (routeHandler *RouteHandler) handleUniqueRequest(w http.ResponseWriter, r *http.Request, key string) {
	valuePostFormValue := r.PostFormValue("value")
	if valuePostFormValue == "" {
		http.Error(w, "Invalid unique value specified", 400)
		return
	}

	valuePostFormTags := r.PostFormValue("tags")
	valuePostFormTags = strings.TrimSpace(valuePostFormTags)
	if valuePostFormTags != "" {
		if validateTags(valuePostFormTags) {
			key += "," + valuePostFormTags
		} else {
			http.Error(w, "Invalid tags specified", 400)
			return
		}
	}

	routeHandler.statter.Unique(key, valuePostFormValue)
}

func validateTags(tagsList string) bool {
	list := strings.Split(tagsList, ",")
	if len(list) == 0 {
		return false
	}

	for _, pair := range list {
		pairItems := strings.Split(pair, "=")
		if len(pairItems) != 2 {
			return false
		} else if len(strings.TrimSpace(pairItems[0])) == 0 {
			return false
		} else if len(strings.TrimSpace(pairItems[1])) == 0 {
			return false
		}
	}

	return true
}
