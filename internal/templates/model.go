package templates

import (
	"github.com/google/uuid"
	"project-a/internal/interfaces"
	"project-a/internal/model"
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
	ContactLists map[*model.List][]*model.Contact
}

type ContactPage struct {
	PageData
	ContactLists map[*model.List][]*model.Contact
	Invitations  []*InvitationTemplate
}

type InvitationTemplate struct {
	InviteUuid uuid.UUID `json:"invite_uuid"`
	IsInviter  bool      `json:"isInviter"`
	Email      string    `json:"email"`
}

type RenderChatArgs struct {
	ur          interfaces.UserRepository
	authService model.Session
	wsUrl       string
}
