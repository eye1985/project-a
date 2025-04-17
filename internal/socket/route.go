package socket

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
	"project-a/internal/user"
)

func RegisterRoutes(m *middleware.Middleware, h *Hub, session shared.Session, ur user.Repository) {
	m.HandleFunc("GET /ws", ServeWs(h, newClient, session, ur))
}
