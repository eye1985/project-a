package health

import "project-a/internal/middleware"

func RegisterRoutes(m *middleware.Middleware, h *Handler) {
	m.HandleFunc("GET /health", h.Health)
}
