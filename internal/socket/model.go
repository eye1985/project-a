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
	FromUuid  uuid.UUID `json:"fromUuid"`
	ToUuid    uuid.UUID `json:"toUuid"`
	Message   string    `json:"message"`
	Event     string    `json:"event"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
