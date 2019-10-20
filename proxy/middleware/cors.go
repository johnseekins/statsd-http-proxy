package middleware

import "net/http"

// validate CORS headers
func ValidateCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		next.ServeHTTP(w, r)
	})
}
