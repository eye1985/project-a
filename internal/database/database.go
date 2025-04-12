package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"path/filepath"
)

func postgresUrl() string {
	pgUrl, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatalf("POSTGRES_URL environment variable not set")
	}

	return pgUrl
}

func migrationPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return "file://" + filepath.Join(dir, "migrations")
}

func Migrate() {

	m, err := migrate.New(migrationPath(), postgresUrl()+"?sslmode=disable")
	if err != nil {
		log.Printf("Error running migration: %v", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Printf("Error running migration: %v", err)
		}

		log.Println("No migrations found")
	}

	log.Println("Migrations applied")
}

func Pool() (*pgxpool.Pool, error) {
	// TODO change ParseConfig
	return pgxpool.New(context.Background(), postgresUrl())
}
