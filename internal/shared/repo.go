package shared

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, email string) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserFromSessionId(ctx context.Context, sessionId string) (*User, error)
	InsertUser(ctx context.Context, username string, email string) (*User, error)
	UpdateUserName(ctx context.Context, username string, userId int64) error
	DeleteUser(ctx context.Context, email string) error
}
