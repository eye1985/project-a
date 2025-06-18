package middleware

import (
	"context"
	"net/http"
	"project-a/internal/consts"
	"project-a/internal/interfaces"
	"project-a/internal/util"
)

func Authenticated(session interfaces.AuthService) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(string(consts.SessionCtxKey))

			if err != nil {
				util.SetFlashCookie(
					w, &util.FlashCookieArgs{
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
				util.SetFlashCookie(
					w, &util.FlashCookieArgs{
						Name:  "flash",
						Value: "unauthorized",
						Path:  "/",
					},
				)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if !session.IsSessionActive(r.Context(), string(cookieValue)) {
				util.SetFlashCookie(
					w, &util.FlashCookieArgs{
						Name:  "flash",
						Value: "unauthorized",
						Path:  "/",
					},
				)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), consts.SessionCtxKey, cookieValue))
			next.ServeHTTP(w, r)
		}
	}
}
