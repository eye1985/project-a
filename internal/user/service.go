package user

import "github.com/jackc/pgx/v5/pgxpool"

type userService struct {
	Repository
}

type Service interface {
	RegisterUser(user *User) error
	GetUsers() ([]*User, error)
}

func (u *userService) RegisterUser(user *User) error {
	err := u.Repository.InsertUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userService) GetUsers() ([]*User, error) {
	return u.Repository.GetUsers()
}

func NewUserService(pool *pgxpool.Pool) Service {
	return &userService{
		Repository: NewUserRepo(pool),
	}
}
