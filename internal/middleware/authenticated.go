package middleware

import (
	"context"
	"net/http"
	"project-a/internal/shared"
)

func Authenticated(session shared.AuthService) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(string(shared.SessionCtxKey))

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

			if !session.IsSessionActive(r.Context(), string(cookieValue)) {
				// TODO add flashcookie
				http.Redirect(w, r, "/?error=unauthorized", http.StatusSeeOther)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), shared.SessionCtxKey, cookieValue))
			next.ServeHTTP(w, r)
		}
	}
}
