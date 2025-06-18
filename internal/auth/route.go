package auth

import (
	"project-a/internal/interfaces"
	"project-a/internal/middleware"
)

func RegisterRoutes(
	m *middleware.Middleware,
	h *Handler,
	as interfaces.AuthService,
	csrf middleware.CSRFHandler,
) {
	m.HandleFunc(
		"POST /createMagicLink",
		h.CreateMagicLinkCode,
		middleware.AllowOnlyPost,
		middleware.AllowOnlyApplicationJson,
		middleware.CSRF(csrf),
	)
	m.HandleFunc(
		"POST /logout", h.Logout,
		middleware.AllowOnlyPost,
		middleware.Authenticated(as),
	)
	m.HandleFunc("GET /signup/{code}", h.RegisterUser, middleware.AllowOnlyGet)
	m.HandleFunc("GET /signin/{code}", h.SignIn, middleware.AllowOnlyGet)
}
