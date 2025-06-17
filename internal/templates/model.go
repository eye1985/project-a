package templates

import (
	"github.com/google/uuid"
	"project-a/internal/models"
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
	ContactLists map[*models.List][]*models.Contact
}

type ContactPage struct {
	PageData
	ContactLists map[*models.List][]*models.Contact
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
