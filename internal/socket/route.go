package socket

import (
	"project-a/internal/contacts"
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(
	m *middleware.Middleware,
	h *Hub,
	as shared.AuthService,
	ur shared.UserRepository,
	cr contacts.Repository,
) {
	m.HandleFunc("GET /ws", ServeWs(h, newClient, as, ur, cr))
}
