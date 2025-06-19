package auth

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/model"
	"time"
)

//go:embed sql/insert_session.sql
var insertSessionSql string

//go:embed sql/get_session_by_user_id.sql
var getSessionByUserIdSql string

//go:embed sql/get_expired_at_from_session_by_session_id.sql
var getExpiredAtFromSessionSql string

//go:embed sql/insert_magic_link.sql
var insertMagicLinkSql string

//go:embed sql/set_active_magic_link.sql
var setActiveMagicLinkSql string

//go:embed sql/delete_session.sql
var deleteSessionSql string

type authRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetSession(ctx context.Context, id int64) (*model.Session, error)
	SetSession(ctx context.Context, args *SetSessionArgs) (*model.Session, error)
	IsSessionActive(ctx context.Context, sessionId string) bool
	CreateMagicLink(ctx context.Context, args *CreateMagicLinkArgs) error
	ActivateNonExpiredMagicLink(ctx context.Context, code string) (*MagicLink, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type SetSessionArgs struct {
	UserID    int64
	SessionID string
	ExpiresAt time.Time
}

func (a *authRepository) SetSession(ctx context.Context, args *SetSessionArgs) (*model.Session, error) {
	userID := args.UserID
	sessionID := args.SessionID
	expiresAt := args.ExpiresAt

	row := a.pool.QueryRow(ctx, insertSessionSql, userID, sessionID, expiresAt)

	session := &model.Session{}
	err := row.Scan(&session.UserId, &session.SessionID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *authRepository) GetSession(ctx context.Context, userId int64) (*model.Session, error) {
	row := a.pool.QueryRow(ctx, getSessionByUserIdSql, userId)

	session := &model.Session{}
	err := row.Scan(&session.UserId, &session.SessionID, &session.ExpiresAt)
	if err != nil {
		return &model.Session{}, err
	}

	return session, nil
}

func (a *authRepository) IsSessionActive(ctx context.Context, sessionId string) bool {
	var expiresAt time.Time
	row := a.pool.QueryRow(ctx, getExpiredAtFromSessionSql, sessionId)
	err := row.Scan(&expiresAt)
	if err != nil {
		return false
	}

	return expiresAt.After(time.Now())
}

type CreateMagicLinkArgs struct {
	email    string
	expiryAt time.Time
	code     string
}

func (a *authRepository) CreateMagicLink(ctx context.Context, args *CreateMagicLinkArgs) error {
	code := args.code
	email := args.email
	expiresAt := args.expiryAt

	_, err := a.pool.Exec(ctx, insertMagicLinkSql, expiresAt, email, code)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) ActivateNonExpiredMagicLink(ctx context.Context, code string) (*MagicLink, error) {
	row := a.pool.QueryRow(ctx, setActiveMagicLinkSql, code, time.Now())
	magicLink := &MagicLink{}
	err := row.Scan(&magicLink.Email)
	if err != nil {
		return nil, err
	}

	return magicLink, nil
}

func (a *authRepository) DeleteSession(ctx context.Context, sessionId string) error {
	conn, err := a.pool.Exec(ctx, deleteSessionSql, sessionId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() == 0 {
		return errSessionNotFound
	}

	return nil
}

func NewRepo(pool *pgxpool.Pool) Repository {
	return &authRepository{pool}
}
