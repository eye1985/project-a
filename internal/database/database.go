package database

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"path/filepath"
)

func Migrate(pgUrl string, rootDir string) {
	mPath := "file://" + filepath.Join(rootDir, "migrations")
	m, err := migrate.New(mPath, pgUrl)
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

	defer func() {
		if sErr, dErr := m.Close(); sErr != nil || dErr != nil {
			log.Fatalf("Error closing migration: %v %v", sErr, dErr)
		}
	}()

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
