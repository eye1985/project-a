package database

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"path/filepath"
)

func migrationPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return "file://" + filepath.Join(dir, "migrations")
}

func Migrate(pgUrl string) {
	m, err := migrate.New(migrationPath(), pgUrl+"?sslmode=disable")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error running migration: %v", err)
		}

		log.Println("No new migrations found")
		return
	}

	log.Println("Migrations applied")
}

func Pool(pgUrl string) (*pgxpool.Pool, error) {
	// TODO change ParseConfig
	return pgxpool.New(context.Background(), pgUrl)
}
