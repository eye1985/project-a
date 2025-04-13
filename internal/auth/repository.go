package auth

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type authRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	GetSession(id int64) (*Session, error)
	SetSession(args *SetSessionArgs) (*Session, error)
	IsSessionActive(sessionId string) bool
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

	row := a.pool.QueryRow(ctx, "INSERT INTO user_sessions (user_id, session_id, expires_at) VALUES ($1, $2, $3) returning user_id, session_id, expires_at", userID, sessionID, expiresAt)

	session := &Session{}
	err := row.Scan(&session.UserId, &session.SessionID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *authRepository) GetSession(userId int64) (*Session, error) {
	ctx := context.Background()
	row := a.pool.QueryRow(ctx, "SELECT user_id, session_id, expires_at FROM user_sessions WHERE user_id = $1", userId)

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
	row := a.pool.QueryRow(ctx, "SELECT expires_at FROM user_sessions WHERE session_id = $1", sessionId)
	err := row.Scan(&expiresAt)
	if err != nil {
		return false
	}

	return expiresAt.After(time.Now())
}

func NewAuthRepo(pool *pgxpool.Pool) Repository {
	return &authRepository{pool}
}
