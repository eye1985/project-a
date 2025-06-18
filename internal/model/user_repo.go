package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        int64     `json:"-"`
	Uuid      uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
