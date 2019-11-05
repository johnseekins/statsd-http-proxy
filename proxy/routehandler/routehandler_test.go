package routehandler

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// create StatsD Client mock
type statsdClientMock struct {
	mock.Mock
}

func (m *statsdClientMock) Open()  {}
func (m *statsdClientMock) Close() {}
func (m *statsdClientMock) Count(key string, value int, sampleRate float32) {
	m.Called(key, value, sampleRate)
}
func (m *statsdClientMock) Timing(key string, time int64, sampleRate float32) {
	m.Called(key, time, sampleRate)
}
func (m *statsdClientMock) Gauge(key string, value int) {
	m.Called(key, value)
}
func (m *statsdClientMock) GaugeShift(key string, value int) {
	m.Called(key, value)
}
func (m *statsdClientMock) Set(key string, value int) {
	m.Called(key, value)
}

func TestHandleHeartbeatRequest(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"GET",
		"http://example.com/heartbeat",
		strings.NewReader(""),
	)

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleHeartbeatRequest(responseWriter, request)

	body, _ := ioutil.ReadAll(responseWriter.Result().Body)

	require.Equal(t, "OK", string(body))
}

func TestHandleCountRequest(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("Count", "someMetricPrefix.a.b.c", 42, float32(0.2)).Return(nil).Once()

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"POST",
		"http://example.com/count/a.b.c",
		strings.NewReader("value=42&sampleRate=0.2"),
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleMetric(responseWriter, request, "count", "a.b.c")

	statsdClient.AssertExpectations(t)
}

func TestHandleSetRequest(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("Set", "someMetricPrefix.a.b.c", 42).Return(nil).Once()

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"POST",
		"http://example.com/set/a.b.c",
		strings.NewReader("value=42"),
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleMetric(responseWriter, request, "set", "a.b.c")

	statsdClient.AssertExpectations(t)
}

func TestHandleTimingRequest(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("Timing", "someMetricPrefix.a.b.c", int64(42000), float32(0.2)).Return(nil).Once()

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"POST",
		"http://example.com/timing/a.b.c",
		strings.NewReader("time=42000&sampleRate=0.2"),
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleMetric(responseWriter, request, "timing", "a.b.c")

	statsdClient.AssertExpectations(t)
}

func TestHandleGaugeRequestWithShift(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("GaugeShift", "someMetricPrefix.a.b.c", 42).Return(nil).Once()

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"POST",
		"http://example.com/gauge/a.b.c",
		strings.NewReader("shift=42"),
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleMetric(responseWriter, request, "gauge", "a.b.c")

	statsdClient.AssertExpectations(t)
}

func TestHandleGaugeRequestWithAbsoluteValue(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("Gauge", "someMetricPrefix.a.b.c", 42).Return(nil).Once()

	// create route handler
	routeHandler := NewRouteHandler(
		statsdClient,
		"someMetricPrefix",
	)

	// mock requerst
	request := httptest.NewRequest(
		"POST",
		"http://example.com/gauge/a.b.c",
		strings.NewReader("value=42"),
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// prepare response
	responseWriter := httptest.NewRecorder()

	// test count handler
	routeHandler.HandleMetric(responseWriter, request, "gauge", "a.b.c")

	statsdClient.AssertExpectations(t)
}
