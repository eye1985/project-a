package users

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

type IUserRepo interface {
	GetUser(email string) (User, error)
	GetUsers() ([]User, error)
	InsertUser(user User) error
	DeleteUser(email string) error
}

type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *UserRepo) GetUser(email string) (User, error) {
	ctx := context.Background()
	var username, userEmail string
	var createdAt time.Time
	row := r.pool.QueryRow(ctx, "select username, email, created_at from users where email=$1", email)
	err := row.Scan(&username, &userEmail, &createdAt)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:  username,
		Email:     userEmail,
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepo) GetUsers() ([]User, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, "select username, email, created_at from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) InsertUser(user User) error {
	ctx := context.Background()
	result, err := r.pool.Exec(ctx, "insert into users(username, email) values($1, $2) ", user.Username, user.Email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No user inserted")
	}

	return nil
}

func (r *UserRepo) DeleteUser(email string) error {
	ctx := context.Background()
	result, err := r.pool.Exec(ctx, "delete from users where email=$1", email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No user deleted")
	}

	return nil
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool}
}

// Convention to ensure UserRepo implements the interface
var _ IUserRepo = (*UserRepo)(nil)
