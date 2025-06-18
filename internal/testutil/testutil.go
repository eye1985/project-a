package testutil

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"project-a/internal/database"
)

func SetupTestContainer(
	ctx context.Context,
) (*postgres.PostgresContainer, *pgxpool.Pool, error) {
	dbName := "testing"
	dbUser := "user"
	dbPass := "password"

	postgresContainer, err := postgres.Run(
		ctx, "postgres:16.9-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return nil, nil, err
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}
	database.Migrate(connStr, "../..")
	//err = postgresContainer.Snapshot(ctx, postgres.WithSnapshotName(snapShotName))
	//if err != nil {
	//	return nil, nil, err
	//}

	pool, err := database.Pool(connStr)
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, pool, nil
}
