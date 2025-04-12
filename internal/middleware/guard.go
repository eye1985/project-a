package middleware

import (
	"net/http"
	"project-a/internal/shared"
)

func Guard(session shared.Session) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sid")
			if err != nil {
				// TODO add flashcookie
				http.Redirect(w, r, "/?error=unauthorized", http.StatusSeeOther)
				return
			}

			if !session.IsSessionActive(cookie.Value) {
				// TODO add flashcookie
				http.Redirect(w, r, "/?error=unauthorized", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
