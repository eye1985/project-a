package auth

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
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

type authRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetSession(id int64) (*Session, error)
	SetSession(args *SetSessionArgs) (*Session, error)
	IsSessionActive(sessionId string) bool
	CreateMagicLink(args *CreateMagicLinkArgs) error
	ActivateMagicLink(code string) (*MagicLink, error)
}

type SetSessionArgs struct {
	userID    int64
	sessionID string
	expiresAt time.Time
}

func (a *authRepository) SetSession(args *SetSessionArgs) (*Session, error) {
	ctx := context.Background()
	userID := args.userID
	sessionID := args.sessionID
	expiresAt := args.expiresAt

	row := a.pool.QueryRow(ctx, insertSessionSql, userID, sessionID, expiresAt)

	session := &Session{}
	err := row.Scan(&session.UserId, &session.SessionID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *authRepository) GetSession(userId int64) (*Session, error) {
	ctx := context.Background()
	row := a.pool.QueryRow(ctx, getSessionByUserIdSql, userId)

	session := &Session{}
	err := row.Scan(&session.UserId, &session.SessionID, &session.ExpiresAt)
	if err != nil {
		return &Session{}, err
	}

	return session, nil
}

func (a *authRepository) IsSessionActive(sessionId string) bool {
	ctx := context.Background()
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

func (a *authRepository) CreateMagicLink(args *CreateMagicLinkArgs) error {
	code := args.code
	email := args.email
	expiresAt := args.expiryAt

	ctx := context.Background()
	_, err := a.pool.Exec(ctx, insertMagicLinkSql, expiresAt, email, code)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) ActivateMagicLink(code string) (*MagicLink, error) {
	ctx := context.Background()
	row := a.pool.QueryRow(ctx, setActiveMagicLinkSql, code, time.Now())
	magicLink := &MagicLink{}
	err := row.Scan(&magicLink.Email)
	if err != nil {
		return nil, err
	}

	return magicLink, nil
}

func NewAuthRepo(pool *pgxpool.Pool) Repository {
	return &authRepository{pool}
}
