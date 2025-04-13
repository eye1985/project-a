package templates

import (
	"project-a/internal/middleware"
	"project-a/internal/shared"
	"project-a/internal/user"
)

type RegisterRoutesArgs struct {
	Middleware  *middleware.Middleware
	WsUrl       string
	Session     shared.Session
	UserService user.Service
}

func RegisterRoutes(args *RegisterRoutesArgs) {
	m := args.Middleware
	wsUrl := args.WsUrl
	session := args.Session
	us := args.UserService

	m.HandleFuncWithMiddleWare("GET /chat", RenderChat(&RenderChatArgs{
		wsUrl:       wsUrl,
		us:          us,
		authService: session,
	}), middleware.Guard(session))
	m.HandleFunc("GET /", RenderRegisterUser)
}
