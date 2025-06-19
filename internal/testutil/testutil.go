package testutil

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Need this for the pgx driver
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"project-a/internal/database"
	"testing"
)

func SetupTestContainer(
	ctx context.Context,
	t *testing.T,
) (*postgres.PostgresContainer, string) {
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
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	database.Migrate(connStr, "../..")

	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)
	err = postgresContainer.Snapshot(ctx)
	require.NoError(t, err)

	return postgresContainer, connStr
}

func CreateTestPoolAndCleanUp(
	t *testing.T,
	ctx context.Context,
	connStr string,
	pgContainer *postgres.PostgresContainer,
) *pgxpool.Pool {
	pool, err := database.Pool(connStr)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			pool.Close()
			err := pgContainer.Restore(ctx)
			require.NoError(t, err)
		},
	)

	return pool
}
