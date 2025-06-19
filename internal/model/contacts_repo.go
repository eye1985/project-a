package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type List struct {
	Id        int64        `json:"-"`
	Uuid      uuid.UUID    `json:"uuid"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at,omitzero"`
	UserId    int64        `json:"user_id"`
}

type Contact struct {
	UserId   int64     `json:"-"`
	UserUuid uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	ListName string    `json:"list_name"`
}

type ContactList struct {
	Id        int64        `json:"-"`
	Uuid      uuid.UUID    `json:"uuid"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt sql.NullTime `json:"updatedAt,omitzero"`
	UserId    int64        `json:"userId"`
}

type InsertedContact struct {
	Id      int64 `json:"-"`
	User1Id int64 `json:"-"`
	User2Id int64 `json:"-"`
}

type Invitation struct {
	Id           int64     `json:"-"`
	Uuid         uuid.UUID `json:"uuid"`
	InviterId    int64     `json:"_"`
	InviteeId    int64     `json:"-"`
	InviterEmail string    `json:"inviterEmail"`
	InviteeEmail string    `json:"inviteeEmail"`
	Accepted     bool      `json:"accepted"`
}

type AcceptedInvite struct {
	Id        int64 `json:"-"`
	InviterId int64 `json:"_"`
	InviteeId int64 `json:"-"`
}
