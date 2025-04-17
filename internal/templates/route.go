package templates

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
	"project-a/internal/user"
)

type RegisterRoutesArgs struct {
	Middleware *middleware.Middleware
	WsUrl      string
	Session    shared.Session
	UserRepo   user.Repository
}

func RegisterRoutes(args *RegisterRoutesArgs) {
	m := args.Middleware
	wsUrl := args.WsUrl
	session := args.Session
	ur := args.UserRepo

	m.HandleFunc("GET /chat", RenderChat(&RenderChatArgs{
		wsUrl:       wsUrl,
		ur:          ur,
		authService: session,
	}), middleware.Guard(session))
	m.HandleFunc("GET /", RenderRegisterUser)
}
