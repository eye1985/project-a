package user

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/interfaces"
	"project-a/internal/model"
)

type userRepository struct {
	pool *pgxpool.Pool
}

//go:embed sql/get_user_by_session_id.sql
var getUserBySessionIdSql string

//go:embed sql/get_user_by_email.sql
var getUserByEmailSql string

//go:embed sql/get_all_users.sql
var getAllUsersSql string

//go:embed sql/insert_user.sql
var insertUserSql string

//go:embed sql/delete_user_by_email.sql
var deleteUserByEmailSql string

//go:embed sql/update_username.sql
var updateUsernameSql string

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	row := r.pool.QueryRow(ctx, getUserByEmailSql, email)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserFromSessionId(ctx context.Context, sessionId string) (*model.User, error) {
	user := &model.User{}
	row := r.pool.QueryRow(ctx, getUserBySessionIdSql, sessionId)
	err := row.Scan(&user.Id, &user.Uuid, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	rows, err := r.pool.Query(ctx, getAllUsersSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.Id, &user.Uuid, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) InsertUser(ctx context.Context, username string, email string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, insertUserSql, username, email)

	u := &model.User{}
	err := row.Scan(
		&u.Id,
		&u.Uuid,
		&u.Email,
		&u.Username,
		&u.CreatedAt,
	)
	if err != nil {
		return &model.User{}, err
	}

	return u, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, email string) error {
	result, err := r.pool.Exec(ctx, deleteUserByEmailSql, email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNoUserDeleted
	}

	return nil
}

func (r *userRepository) UpdateUserName(ctx context.Context, newUsername string, userId int64) error {
	result, err := r.pool.Exec(ctx, updateUsernameSql, newUsername, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNoUsernameUpdated
	}

	return nil
}

func NewUserRepo(pool *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{pool}
}
