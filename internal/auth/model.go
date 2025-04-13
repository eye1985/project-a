package auth

import "time"

type Session struct {
	UserId    int64
	SessionID string
	ExpiresAt time.Time
}

type MagicLink struct {
	Email string
}
