package auth

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, as shared.AuthService) {
	m.HandleFunc("POST /login", h.Login, middleware.AllowOnlyPost)
	m.HandleFunc("POST /logout", h.Logout,
		middleware.AllowOnlyPost,
		middleware.Guard(as),
	)
	m.HandleFunc("GET /signup/{code}", h.RegisterUser)
}
