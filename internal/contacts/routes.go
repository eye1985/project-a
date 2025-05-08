package contacts

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, authService shared.AuthService) {
	m.HandleFunc("POST /contactlist", h.CreateUserList, middleware.Authenticated(authService))
	m.HandleFunc(
		"POST /contact",
		h.CreateInvitation,
		middleware.Authenticated(authService),
		middleware.AllowOnlyPost,
		middleware.AllowOnlyApplicationJson,
	)
	m.HandleFunc(
		"PATCH /invite",
		h.AcceptInvite,
		middleware.Authenticated(authService),
		middleware.AllowOnlyPatch,
		middleware.AllowOnlyApplicationJson,
	)
}
