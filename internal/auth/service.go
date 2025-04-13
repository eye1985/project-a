package auth

import (
	"encoding/base64"
	"github.com/gorilla/securecookie"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

type authService struct {
	Repository
	SecureCookie *securecookie.SecureCookie
}

type Service interface {
	CreateOrGetSession(userId int64) (*Session, error)
	IsSessionActive(sessionId string) bool
	SignCookie(cookieName string, value []byte) (string, error)
	VerifyCookie(cookie *http.Cookie) ([]byte, error)
}

func (a *authService) CreateOrGetSession(userId int64) (*Session, error) {
	s, err := a.Repository.GetSession(userId)
	log.Printf("create or get session:  %s", s.SessionID)
	log.Printf("error %s", err)
	if err != nil {
		// TODO check if its actually an not found error
		sessionID, err := createSessionID()
		if err != nil {
			return nil, err
		}
		// No session, register a new session
		const thirtyDays = 30 * 24 * time.Hour
		ns, err := a.SetSession(&SetSessionArgs{
			userID:    userId,
			sessionID: sessionID,
			expiresAt: time.Now().Add(thirtyDays),
		})

		if err != nil {
			return nil, err
		}

		return ns, nil
	}

	return s, nil
}

func (a *authService) IsSessionActive(sessionId string) bool {
	return a.Repository.IsSessionActive(sessionId)
}

func (a *authService) SignCookie(cookieName string, value []byte) (string, error) {
	encoded, err := a.SecureCookie.Encode(cookieName, value)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return encoded, nil
}

func (a *authService) VerifyCookie(cookie *http.Cookie) ([]byte, error) {
	var decoded []byte
	err := a.SecureCookie.Decode(cookie.Name, cookie.Value, &decoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func NewAuthService(pool *pgxpool.Pool, hashKey string, blockKey string) Service {
	hk, err := base64.StdEncoding.DecodeString(hashKey)
	if err != nil {
		log.Fatalf("Failed to decode hash key: %s", err)
	}

	bk, err := base64.StdEncoding.DecodeString(blockKey)
	if err != nil {
		log.Fatalf("Failed to decode block key: %s", err)
	}

	return &authService{
		Repository:   NewAuthRepo(pool),
		SecureCookie: securecookie.New(hk, bk),
	}
}
