package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type authService struct {
	Repository
}

type Service interface {
	CreateOrGetSession(userId int64) (*Session, error)
	IsSessionActive(sessionId string) bool
}

func (a *authService) CreateOrGetSession(userId int64) (*Session, error) {
	s, err := a.GetSession(userId)
	if err != nil {
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

func NewAuthService(pool *pgxpool.Pool) Service {
	return &authService{
		Repository: NewAuthRepo(pool),
	}
}
