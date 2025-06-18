package socket

import (
	"project-a/internal/interfaces"
	"project-a/internal/middleware"
)

func RegisterRoutes(
	m *middleware.Middleware,
	h *Hub,
	as interfaces.AuthService,
	ur interfaces.UserRepository,
	cr interfaces.ContactsRepository,
	origin string,
) {
	m.HandleFunc("GET /ws", ServeWs(h, newClient, as, ur, cr, origin))
}
