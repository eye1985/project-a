package auth

import (
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"project-a/internal/shared"
	"time"
)

type authService struct {
	Repository
	SecureCookie *securecookie.SecureCookie
}

func (a *authService) CreateOrGetSession(ctx context.Context, userId int64) (*shared.Session, error) {
	s, err := a.Repository.GetSession(ctx, userId)

	if err != nil {
		// TODO check if its actually an not found error
		sessionID, err := createSessionID()
		if err != nil {
			return nil, err
		}
		// No session, register a new session
		const thirtyDays = 30 * 24 * time.Hour
		ns, err := a.SetSession(ctx, &SetSessionArgs{
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

// IsSessionActive Used in guard
func (a *authService) IsSessionActive(ctx context.Context, sessionId string) bool {
	return a.Repository.IsSessionActive(ctx, sessionId)
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

func (a *authService) CreateMagicLink(ctx context.Context, email string) (string, error) {
	tenMinFromNow := time.Now().Add(10 * time.Minute)
	u := uuid.New()
	encoding := base64.URLEncoding.WithPadding(base64.NoPadding)
	encoded := encoding.EncodeToString(u[:])

	err := a.Repository.CreateMagicLink(ctx, &CreateMagicLinkArgs{
		code:     encoded,
		expiryAt: tenMinFromNow,
		email:    email,
	})

	if err != nil {
		return "", err
	}

	return encoded, nil
}

func NewService(repo Repository, hashKey string, blockKey string) shared.AuthService {
	hk, err := base64.StdEncoding.DecodeString(hashKey)
	if err != nil {
		log.Fatalf("Failed to decode hash key: %s", err)
	}

	bk, err := base64.StdEncoding.DecodeString(blockKey)
	if err != nil {
		log.Fatalf("Failed to decode block key: %s", err)
	}

	return &authService{
		Repository:   repo,
		SecureCookie: securecookie.New(hk, bk),
	}
}
