package shared

type UserRepository interface {
	GetUser(email string) (*User, error)
	GetUsers() ([]*User, error)
	GetUserFromSessionId(sessionId string) (*User, error)
	InsertUser(username string, email string) (*User, error)
	UpdateUserName(username string, userId int64) error
	DeleteUser(email string) error
}
