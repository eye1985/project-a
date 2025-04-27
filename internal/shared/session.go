package shared

import (
	"context"
	"net/http"
)

type AuthService interface {
	CreateOrGetSession(ctx context.Context, userId int64) (*Session, error)
	IsSessionActive(ctx context.Context, sessionId string) bool
	SignCookie(cookieName string, value []byte) (string, error)
	VerifyCookie(cookie *http.Cookie) ([]byte, error)
	CreateMagicLink(ctx context.Context, email string) (string, error)
}

type ctxKey string

const SessionCtxKey = ctxKey("sid")
