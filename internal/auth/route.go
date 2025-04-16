package auth

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, session shared.Session) {
	m.HandleFunc("POST /login", h.Login)
	m.HandleFunc("POST /logout", h.Logout, middleware.Guard(session))
	m.HandleFunc("GET /signup/{code}", h.RegisterUser)
}
