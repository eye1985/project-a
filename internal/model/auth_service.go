package model

import "time"

type Session struct {
	UserId    int64
	SessionID string
	ExpiresAt time.Time
}
