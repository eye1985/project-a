package templates

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

type RegisterRoutesArgs struct {
	Middleware *middleware.Middleware
	WsUrl      string
	Session    shared.Session
	UserRepo   shared.UserRepository
}

func RegisterRoutes(m *middleware.Middleware, h *Handler, authService shared.Session) {
	m.HandleFunc("GET /", h.RenderRegisterUser)
	m.HandleFunc("GET /chat", h.RenderChat, middleware.Guard(authService))
	m.HandleFunc("GET /profile", h.RenderProfile, middleware.Guard(authService))
	m.HandleFunc("GET /userlists", h.RenderUserList, middleware.Guard(authService))
}
