package routehandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type CountRequest struct {
	Value    int    `json:"value,omitempty"`
	Tags string `json:"tags,omitempty"`
}

const maxBodySize = 10 * 1024 * 1024

func (routeHandler *RouteHandler) handleCountRequest(w http.ResponseWriter, r *http.Request, key string) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Unsupported content type", 400)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	r.Body.Close()

	var req CountRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	reqTags := strings.TrimSpace(req.Tags)
	if validateTags(reqTags) {
		key += "," + reqTags
	}

	routeHandler.statter.Count(key, req.Value)
}

type IncrRequest struct {
	Tags string `json:"tags,omitempty"`
}

func (routeHandler *RouteHandler) handleIncrementRequest(w http.ResponseWriter, r *http.Request, key string) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Unsupported content type", 400)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	r.Body.Close()

	var req IncrRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	reqTags := strings.TrimSpace(req.Tags)
	if validateTags(reqTags) {
		key += "," + reqTags
	}

	routeHandler.statter.Increment(key)
}

type GaugeRequest struct {
	Value int    `json:"value,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

func (routeHandler *RouteHandler) handleGaugeRequest(w http.ResponseWriter, r *http.Request, key string) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Unsupported content type", 400)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	r.Body.Close()

	var req GaugeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	reqTags := strings.TrimSpace(req.Tags)
	if validateTags(reqTags) {
		key += "," + reqTags
	}

	routeHandler.statter.Gauge(key, req.Value)
}

type TimingRequest struct {
	Value int    `json:"value,omitempty"`
	Tags     string `json:"tags,omitempty"`
}

func (routeHandler *RouteHandler) handleTimingRequest(w http.ResponseWriter, r *http.Request, key string) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Unsupported content type", 400)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	r.Body.Close()

	var req TimingRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	reqTags := strings.TrimSpace(req.Tags)
	if validateTags(reqTags) {
		key += "," + reqTags
	}

	routeHandler.statter.Timing(key, req.Value)
}

type UniqueRequest struct {
	Value string `json:"value,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

func (routeHandler *RouteHandler) handleUniqueRequest(w http.ResponseWriter, r *http.Request, key string) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Unsupported content type", 400)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	r.Body.Close()

	var req UniqueRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	reqTags := strings.TrimSpace(req.Tags)
	if validateTags(reqTags) {
		key += "," + reqTags
	}

	routeHandler.statter.Unique(key, req.Value)
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
