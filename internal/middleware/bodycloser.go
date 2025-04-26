package middleware

import (
	"log"
	"net/http"
)

func BodyCloser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Printf("failed to close request body: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
