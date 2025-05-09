package middleware

import "net/http"

func NoCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store") // Dont store response at all
		w.Header().Set("Pragma", "no-cache")        // Legacy http/1.0 header
		w.Header().Set("Expires", "0")              // Legacy header, resource already expired

		next.ServeHTTP(w, r)
	}
}
