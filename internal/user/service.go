package user

import "github.com/jackc/pgx/v5/pgxpool"

type userService struct {
	Repository
}

type Service interface {
	RegisterUser(user *User) (*User, error)
	GetUsers() ([]*User, error)
	GetUser(email string) (*User, error)
}

func (u *userService) RegisterUser(user *User) (*User, error) {
	newUser, err := u.Repository.InsertUser(user)

	if err != nil {
		return &User{}, err
	}

	return newUser, nil
}

func (u *userService) GetUsers() ([]*User, error) {
	return u.Repository.GetUsers()
}

func (u *userService) GetUser(email string) (*User, error) {
	return u.Repository.GetUser(email)
}

func NewUserService(pool *pgxpool.Pool) Service {
	return &userService{
		Repository: NewUserRepo(pool),
	}
}
