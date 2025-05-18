package email

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"net"
	"project-a/internal/database"
	"project-a/migrations"
	"testing"
)

func TestWithPostgresSQL(t *testing.T) {
	ctx := context.Background()

	dbName := "testing"
	dbUser := "user"
	dbPass := "password"

	postgresContainer, err := postgres.Run(
		ctx, "postgres:16.8-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)

	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	// Run any migrations on the database
	_, _, err = postgresContainer.Exec(
		ctx,
		[]string{
			"psql",
			"-U",
			dbUser,
			"-d",
			dbName,
			"-c",
			migrations.InitMigrationCreate,
		},
	)
	require.NoError(t, err)

	err = postgresContainer.Snapshot(ctx, postgres.WithSnapshotName("sentEmail-snapshot"))
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := database.Pool(connStr)
	require.NoError(t, err)

	t.Run(
		"Should add one email row", func(t *testing.T) {
			t.Cleanup(
				func() {
					// Restore the database to the snapshot state
					err = postgresContainer.Restore(ctx)
					require.NoError(t, err)
				},
			)

			repo := NewRepo(pool)
			err := repo.AddSentEmail(ctx, "test@test.com", net.ParseIP("127.0.0.1").To4())
			require.NoError(t, err)

			emails, err := repo.GetSentEmails(ctx)
			require.NoError(t, err)
			require.Len(t, emails, 1)
		},
	)
}
