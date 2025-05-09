package cookieutil

import (
	"net/http"
)

type FlashCookieArgs struct {
	Name  string
	Value string
	Path  string
}

func SetFlashCookie(w http.ResponseWriter, args *FlashCookieArgs) {
	http.SetCookie(
		w, &http.Cookie{
			Name:     args.Name,
			Value:    args.Value,
			Path:     args.Path,
			MaxAge:   60,
			Secure:   true,
			HttpOnly: false,
			SameSite: http.SameSiteStrictMode,
		},
	)
}
