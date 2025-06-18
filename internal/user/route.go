package user

import (
	"project-a/internal/interfaces"
	"project-a/internal/middleware"
)

func RegisterRoutes(m *middleware.Middleware, h *Handler, as interfaces.AuthService) {
	m.HandleFunc("GET /users", h.GetUsers)
	//m.HandleFunc("POST /users", h.RegisterUser,
	//	middleware.AllowOnlyPost,
	//	middleware.AllowOnlyApplicationJson,
	//)
	m.HandleFunc("PATCH /user", h.UpdateUserName,
		middleware.AllowOnlyPatch,
		middleware.AllowOnlyApplicationJson,
		middleware.Authenticated(as),
	)
}
