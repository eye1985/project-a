package socket

import (
	"project-a/internal/interfaces"
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(
	m *middleware.Middleware,
	h *Hub,
	as shared.AuthService,
	ur shared.UserRepository,
	cr interfaces.ContactsRepository,
	origin string,
) {
	m.HandleFunc("GET /ws", ServeWs(h, newClient, as, ur, cr, origin))
}
