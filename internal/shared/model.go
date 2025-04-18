package shared

import "time"

type User struct {
	Id        int64     `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
