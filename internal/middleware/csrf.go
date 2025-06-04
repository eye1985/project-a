package middleware

import (
	"net/http"
)

func CSRF(csrf CSRFHandler) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		h := csrf(next)
		return func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
}
