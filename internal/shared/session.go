package shared

type Session interface {
	IsSessionActive(sessionId string) bool
}
