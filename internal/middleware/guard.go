package middleware

import (
	"context"
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

			cookieValue, err := session.VerifyCookie(cookie)
			if err != nil {
				http.Redirect(w, r, "/?error=unauthorized", http.StatusSeeOther)
				return
			}

			if !session.IsSessionActive(string(cookieValue)) {
				// TODO add flashcookie
				http.Redirect(w, r, "/?error=unauthorized", http.StatusSeeOther)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), shared.SessionCtxKey, cookieValue))
			next.ServeHTTP(w, r)
		}
	}
}
