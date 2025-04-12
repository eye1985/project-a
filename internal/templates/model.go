package templates

import "github.com/jackc/pgx/v5/pgxpool"

type Person struct {
	Username string
	Email    string
}

type PageData struct {
	Person []Person
	WsUrl  string
}

type RenderChatArgs struct {
	Pool  *pgxpool.Pool
	WsUrl string
}
