package routehandler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
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
func (m *statsdClientMock) Timing(key string, time int64, sampleRate float32) {}
func (m *statsdClientMock) Gauge(key string, value int)                       {}
func (m *statsdClientMock) GaugeShift(key string, value int)                  {}
func (m *statsdClientMock) Set(key string, value int)                         {}

func TestHandleCountRequest(t *testing.T) {
	// create statsd client mock
	statsdClient := new(statsdClientMock)

	statsdClient.On("Count", "a.b.c", 42, float32(0.2)).Return(nil).Once()

	// create route handler
	routeHandler := RouteHandler{
		statsdClient,
		"someMetricPrefix",
	}

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
	routeHandler.handleCountRequest(responseWriter, request, "a.b.c")

	statsdClient.AssertExpectations(t)
}
