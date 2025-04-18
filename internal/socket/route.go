package socket

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Hub, session shared.Session, ur shared.UserRepository) {
	m.HandleFunc("GET /ws", ServeWs(h, newClient, session, ur))
}
