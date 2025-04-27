package userlists

import "time"

type UserList struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    int64     `json:"user_id"`
}

type CreateUserListBody struct {
	Name   string `json:"name"`
	UserId int64  `json:"user_id"`
}
