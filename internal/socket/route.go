package socket

import "project-a/internal/middleware"

func RegisterRoutes(m *middleware.Middleware, h *Hub) {
	m.HandleFunc("GET /ws", ServeWs(h, NewClient))
}
