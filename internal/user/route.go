package user

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, authService shared.Session) {
	m.HandleFunc("GET /users", h.GetUsers)
	m.HandleFunc("POST /users", h.RegisterUser)
	m.HandleFunc("PATCH /user", h.UpdateUserName, middleware.Guard(authService))
}
