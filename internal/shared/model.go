package shared

import "time"

type Session struct {
	UserId    int64
	SessionID string
	ExpiresAt time.Time
}
type User struct {
	Id        int64     `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
