package socket

import (
	"github.com/google/uuid"
	"time"
)

type MessageIn struct {
	ToUuid uuid.UUID `json:"toUuid"`
	Msg    string    `json:"msg"`
}

type MessageJSON struct {
	Uuid      uuid.UUID `json:"uuid"`
	Message   string    `json:"message"`
	Event     string    `json:"event"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
