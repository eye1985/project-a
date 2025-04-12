package user

import (
	"project-a/internal/middleware"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler) {
	m.HandleFunc("GET /users", h.GetUsers)
	m.HandleFunc("POST /users", h.RegisterUser)
}
