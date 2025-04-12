package auth

import "project-a/internal/middleware"

func RegisterRoutes(m *middleware.Middleware, h *Handler) {
	m.HandleFunc("POST /login", h.Login)
}
