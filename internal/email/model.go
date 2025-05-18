package email

import (
	"net"
	"time"
)

type SentEmail struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Ip        net.IP    `json:"ip"`
	IsSignUp  bool      `json:"isSignUp"`
}
