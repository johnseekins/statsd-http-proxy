package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCORSWithoutOriginHeaderInRequest(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	handlerWithJWTValidation := ValidateCORS(nextHandler)

	request := httptest.NewRequest("GET", "http://testing", nil)
	responseWriter := httptest.NewRecorder()

	handlerWithJWTValidation.ServeHTTP(responseWriter, request)

	response := responseWriter.Result()

	require := require.New(t)

	require.Empty(response.Header.Get("Access-Control-Allow-Origin"))
}

func TestValidateCORSWithPreflightRequest(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	handlerWithJWTValidation := ValidateCORS(nextHandler)

	request := httptest.NewRequest("OPTIONS", "http://testing", nil)
	responseWriter := httptest.NewRecorder()

	origin := "example.com"
	request.Header.Add("Origin", origin)
	request.Header.Add("Access-Control-Request-Method", "GET")
	request.Header.Add("Access-Control-Request-Headers", "GET")

	handlerWithJWTValidation.ServeHTTP(responseWriter, request)

	response := responseWriter.Result()

	require := require.New(t)

	require.Equal(http.StatusNoContent, response.StatusCode)
	require.Equal(origin, response.Header.Get("Access-Control-Allow-Origin"))

	require.NotEmpty(response.Header.Get("Access-Control-Allow-Methods"))

	require.NotEmpty(response.Header.Get("Access-Control-Allow-Headers"))
}
