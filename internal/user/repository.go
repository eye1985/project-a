package user

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

//go:embed sql/get_user_by_session_id.sql
var getUserBySessionId string

//go:embed sql/get_user_by_email.sql
var getUserByEmail string

//go:embed sql/get_all_users.sql
var getAllUsers string

//go:embed sql/insert_user.sql
var insertUser string

//go:embed sql/delete_user_by_email.sql
var deleteUserByEmail string

type Repository interface {
	GetUser(email string) (*User, error)
	GetUsers() ([]*User, error)
	GetUserFromSessionId(sessionId string) (*User, error)
	InsertUser(user *User) (*User, error)
	DeleteUser(email string) error
}

func (r *userRepository) GetUser(email string) (*User, error) {
	ctx := context.Background()
	user := &User{}

	row := r.pool.QueryRow(ctx, getUserByEmail, email)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserFromSessionId(sessionId string) (*User, error) {
	ctx := context.Background()
	user := &User{}
	row := r.pool.QueryRow(ctx, getUserBySessionId, sessionId)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUsers() ([]*User, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) InsertUser(user *User) (*User, error) {
	ctx := context.Background()
	row := r.pool.QueryRow(ctx, insertUser, user.Username, user.Email)

	u := &User{}
	err := row.Scan(&u.Id, &u.Email, &u.Username, &u.CreatedAt)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (r *userRepository) DeleteUser(email string) error {
	ctx := context.Background()
	result, err := r.pool.Exec(ctx, deleteUserByEmail, email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No user deleted")
	}

	return nil
}

func NewUserRepo(pool *pgxpool.Pool) Repository {
	return &userRepository{pool}
}
