package contacts

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, authService shared.AuthService) {
	m.HandleFunc("POST /contactlist", h.CreateUserList, middleware.Guard(authService))
	m.HandleFunc(
		"POST /contact",
		h.CreateContact,
		middleware.Guard(authService),
		middleware.AllowOnlyPost,
		middleware.AllowOnlyApplicationJson,
	)
	m.HandleFunc(
		"PATCH /contact",
		h.AcceptInvite,
		middleware.Guard(authService),
		middleware.AllowOnlyPatch,
		middleware.AllowOnlyApplicationJson,
	)
}
