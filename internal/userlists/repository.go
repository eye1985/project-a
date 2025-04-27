package userlists

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type userListRepository struct {
	pool *pgxpool.Pool
}

//go:embed sql/create_user_list.sql
var insertUserListSql string

//go:embed sql/update_user_list.sql
var updateUserListSql string

//go:embed sql/delete_user_list.sql
var deleteUserListSql string

//go:embed sql/get_all_user_lists.sql
var getAllUserListsSql string

//go:embed sql/get_user_list.sql
var getUserListSql string

type UserListRepository interface {
	GetUserLists(ctx context.Context, userId int64) ([]*UserList, error)
	GetUserList(ctx context.Context, userListId int64) (*UserList, error)
	CreateUserList(ctx context.Context, name string, userId int64) error
	UpdateUserList(ctx context.Context, name string, userListId int64) error
	DeleteUserList(ctx context.Context, userListId int64) error
}

func (ulr *userListRepository) CreateUserList(ctx context.Context, name string, userId int64) error {
	conn, err := ulr.pool.Exec(ctx, insertUserListSql, name, userId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return userListNotCreated
	}

	return nil
}

func (ulr *userListRepository) UpdateUserList(ctx context.Context, name string, userListId int64) error {
	conn, err := ulr.pool.Exec(ctx, updateUserListSql, name, time.Now(), userListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return userListNotUpdated
	}

	return nil
}

func (ulr *userListRepository) DeleteUserList(ctx context.Context, userListId int64) error {
	conn, err := ulr.pool.Exec(ctx, deleteUserListSql, userListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return userListNotDeleted
	}

	return nil
}

func (ulr *userListRepository) GetUserLists(ctx context.Context, userId int64) ([]*UserList, error) {
	rows, err := ulr.pool.Query(ctx, getAllUserListsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLists []*UserList
	for rows.Next() {
		var userList UserList
		err = rows.Scan(
			&userList.Id,
			&userList.Name,
			&userList.CreatedAt,
			&userList.UpdatedAt,
			&userList.UserId,
		)
		if err != nil {
			return nil, err
		}
		userLists = append(userLists, &userList)
	}

	return userLists, nil
}

func (ulr *userListRepository) GetUserList(ctx context.Context, userListId int64) (*UserList, error) {
	row := ulr.pool.QueryRow(ctx, getUserListSql, userListId)
	var userList UserList
	err := row.Scan(
		&userList.Id,
		&userList.Name,
		&userList.CreatedAt,
		&userList.UpdatedAt,
		&userList.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &userList, nil
}

func NewUserListsRepository(pool *pgxpool.Pool) UserListRepository {
	return &userListRepository{
		pool: pool,
	}
}
