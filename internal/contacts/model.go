package contacts

import (
	"github.com/google/uuid"
	"time"
)

type List struct {
	Id        int64      `json:"-"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	UserId    int64      `json:"user_id"`
}

type Contact struct {
	Id              int64  `json:"-"`
	Uuid            string `json:"uuid"`
	InviterId       int64  `json:"-"`
	InviterEmail    string `json:"inviterEmail"`
	InviterUsername string `json:"inviterUsername"`
	InviteeId       int64  `json:"-"`
	InviteeEmail    string `json:"inviteeEmail"`
	InviteeUsername string `json:"inviteeUsername"`
	HasAccepted     bool   `json:"hasAccepted"`
	IsInviter       bool   `json:"isInviter"`
}

type CreateUserListBody struct {
	Name   string `json:"name"`
	UserId int64  `json:"user_id"`
}

type CreateContactBody struct {
	Email         string `json:"email"`
	ContactListId int64  `json:"contactListId"`
}

type AcceptInviteBody struct {
	Uuid uuid.UUID `json:"uuid"`
}
