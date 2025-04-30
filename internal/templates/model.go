package templates

import (
	"project-a/internal/contacts"
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

type ContactListPage struct {
	Title        string
	Css          string
	Username     string
	ContactLists map[*contacts.List][]*contacts.Contact
}

type RenderChatArgs struct {
	ur          shared.UserRepository
	authService shared.Session
	wsUrl       string
}
