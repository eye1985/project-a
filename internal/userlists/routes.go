package userlists

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, authService shared.Session) {
	m.HandleFunc("POST /userlist", h.CreateUserList, middleware.Guard(authService))
}
