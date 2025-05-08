package auth

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, as shared.AuthService) {
	m.HandleFunc("POST /createMagicLink", h.CreateMagicLinkCode, middleware.AllowOnlyPost)
	m.HandleFunc(
		"POST /logout", h.Logout,
		middleware.AllowOnlyPost,
		middleware.Authenticated(as),
	)
	m.HandleFunc("GET /signup/{code}", h.RegisterUser, middleware.AllowOnlyGet)
	m.HandleFunc("GET /signin/{code}", h.SignIn, middleware.AllowOnlyGet)
}
