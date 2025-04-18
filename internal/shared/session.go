package shared

import "net/http"

// TODO refactor this, we dont need double of the same interface as auth
type Session interface {
	IsSessionActive(sessionId string) bool
	SignCookie(cookieName string, value []byte) (string, error)
	VerifyCookie(cookie *http.Cookie) ([]byte, error)
}

type ctxKey string

const SessionCtxKey = ctxKey("sid")
