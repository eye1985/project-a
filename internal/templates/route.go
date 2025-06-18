package templates

import (
	"project-a/internal/interfaces"
	"project-a/internal/middleware"
	"project-a/internal/model"
)

type RegisterRoutesArgs struct {
	Middleware *middleware.Middleware
	WsUrl      string
	Session    model.Session
	UserRepo   interfaces.UserRepository
}

func RegisterRoutes(
	m *middleware.Middleware,
	h *Handler, authService interfaces.AuthService,
	csrf middleware.CSRFHandler,
) {
	m.HandleFunc("GET /", h.RenderRegisterUser, middleware.CSRF(csrf))
	m.HandleFunc("GET /profile", h.RenderProfile, middleware.Authenticated(authService))
	m.HandleFunc("GET /chat", h.RenderChat, middleware.Authenticated(authService))
	m.HandleFunc("GET /contacts", h.RenderContacts, middleware.Authenticated(authService))
}
