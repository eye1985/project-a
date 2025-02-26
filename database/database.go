package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func postgresUrl() string {
	pgUrl, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		pgUrl = "postgres://admin:root@postgres:5432/project_a"
	}

	return pgUrl
}

func Pool() (*pgxpool.Pool, error) {
	// TODO change ParseConfig
	return pgxpool.New(context.Background(), postgresUrl())
}
