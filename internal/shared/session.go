package shared

import "net/http"

type Session interface {
	IsSessionActive(sessionId string) bool
	SignCookie(cookieName string, value []byte) (string, error)
	VerifyCookie(cookie *http.Cookie) ([]byte, error)
}
