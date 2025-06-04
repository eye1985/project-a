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
	Uuid     uuid.UUID
	Title    string
	Css      []string
	Js       string
}

type CreateMagicLinkPage struct {
	PageData
	CSRF string
}
type ChatPage struct {
	PageData
	ContactLists map[*contacts.List][]*contacts.Contact
}

type ContactPage struct {
	PageData
	ContactLists map[*contacts.List][]*contacts.Contact
	Invitations  []*InvitationTemplate
}

type InvitationTemplate struct {
	InviteUuid uuid.UUID `json:"invite_uuid"`
	IsInviter  bool      `json:"isInviter"`
	Email      string    `json:"email"`
}

type RenderChatArgs struct {
	ur          shared.UserRepository
	authService shared.Session
	wsUrl       string
}
