package templates

import (
	"project-a/internal/shared"
	"project-a/internal/user"
)

type Person struct {
	Username string
	Email    string
}

type PageData struct {
	WsUrl    string
	Username string
}

type RenderChatArgs struct {
	us          user.Service
	authService shared.Session
	wsUrl       string
}
