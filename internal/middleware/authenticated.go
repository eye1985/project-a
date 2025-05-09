package middleware

import (
	"context"
	"net/http"
	"project-a/internal/cookieutil"
	"project-a/internal/shared"
)

func Authenticated(session shared.AuthService) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(string(shared.SessionCtxKey))

			if err != nil {
				cookieutil.SetFlashCookie(
					w, &cookieutil.FlashCookieArgs{
						Name:  "flash",
						Value: "unauthorized",
						Path:  "/",
					},
				)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			cookieValue, err := session.VerifyCookie(cookie)
			if err != nil {
				cookieutil.SetFlashCookie(
					w, &cookieutil.FlashCookieArgs{
						Name:  "flash",
						Value: "unauthorized",
						Path:  "/",
					},
				)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if !session.IsSessionActive(r.Context(), string(cookieValue)) {
				cookieutil.SetFlashCookie(
					w, &cookieutil.FlashCookieArgs{
						Name:  "flash",
						Value: "unauthorized",
						Path:  "/",
					},
				)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), shared.SessionCtxKey, cookieValue))
			next.ServeHTTP(w, r)
		}
	}
}
