package interfaces

import (
	"context"
	"project-a/internal/model"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUserFromSessionId(ctx context.Context, sessionId string) (*model.User, error)
	InsertUser(ctx context.Context, username string, email string) (*model.User, error)
	UpdateUserName(ctx context.Context, username string, userId int64) error
	DeleteUser(ctx context.Context, email string) error
}
