package templates

import (
	"github.com/google/uuid"
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

type InvitationTemplate struct {
	InviteUuid uuid.UUID `json:"invite_uuid"`
	IsInviter  bool      `json:"isInviter"`
	Email      string    `json:"email"`
}

type ContactListPage struct {
	Title        string
	Css          string
	Username     string
	Uuid         uuid.UUID
	ContactLists map[*contacts.List][]*contacts.Contact
	Invitations  []*InvitationTemplate
	WsUrl        string
}

type RenderChatArgs struct {
	ur          shared.UserRepository
	authService shared.Session
	wsUrl       string
}
