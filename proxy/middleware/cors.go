package middleware

import "net/http"

const jwtHeaderName = "X-JWT-Token"

// validate CORS headers
func validateCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		next.ServeHTTP(w, r)
	})
}
