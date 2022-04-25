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
	SampleRate float64 `json:"sampleRate"`
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

    key += processTags(req.Tags)	

	var sampleRate float64 = 1
	if req.SampleRate != 0 {
		sampleRate = float64(req.SampleRate)
	}

	routeHandler.statsdClient.Count(key, req.Value, float32(sampleRate))
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

    key += processTags(req.Tags)	

	routeHandler.statsdClient.Gauge(key, req.Value)
}

type TimingRequest struct {
	Value int64    `json:"value,omitempty"`
	Tags     string `json:"tags,omitempty"`
	SampleRate float64 `json:"sampleRate"`
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

    key += processTags(req.Tags)	

	var sampleRate float64 = 1
	if req.SampleRate != 0 {
		sampleRate = float64(req.SampleRate)
	}


	routeHandler.statsdClient.Timing(key, req.Value, float32(sampleRate))
}

type SetRequest struct {
	Value int `json:"value,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

func (routeHandler *RouteHandler) handleSetRequest(w http.ResponseWriter, r *http.Request, key string) {
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

	var req SetRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

    key += processTags(req.Tags)	

	routeHandler.statsdClient.Set(key, req.Value)
}

func processTags(tagsList string) string {
	list := strings.Split(strings.TrimSpace(tagsList), ",")
	if len(list) == 0 {
		return ""
	}

	for _, pair := range list {
		pairItems := strings.Split(pair, "=")
		if len(pairItems) != 2 {
			return ""
		} else if len(strings.TrimSpace(pairItems[0])) == 0 {
			return ""
		} else if len(strings.TrimSpace(pairItems[1])) == 0 {
			return ""
		}
	}

	return "," + tagsList
}
