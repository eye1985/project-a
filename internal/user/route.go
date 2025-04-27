package user

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, as shared.AuthService) {
	m.HandleFunc("GET /users", h.GetUsers)
	m.HandleFunc("POST /users", h.RegisterUser,
		middleware.AllowOnlyPost,
		middleware.AllowOnlyApplicationJson,
	)
	m.HandleFunc("PATCH /user", h.UpdateUserName,
		middleware.AllowOnlyPatch,
		middleware.AllowOnlyApplicationJson,
		middleware.Guard(as),
	)
}
