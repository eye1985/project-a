package templates

import (
	"project-a/internal/shared"
)

type Person struct {
	Username string
	Email    string
}

type PageData struct {
	WsUrl    string
	Username string
	Title    string
	Css      string
}

type RenderChatArgs struct {
	ur          shared.UserRepository
	authService shared.Session
	wsUrl       string
}
