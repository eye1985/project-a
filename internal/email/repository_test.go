package email

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"net"
	"os"
	"project-a/internal/database"
	"testing"
)

func TestWithPostgresSQL(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

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

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	root := os.Getenv("PROJECT_ROOT")

	database.Migrate(connStr, root)
	err = postgresContainer.Snapshot(ctx, postgres.WithSnapshotName("sentEmail-snapshot"))
	require.NoError(t, err)

	pool, err := database.Pool(connStr)
	require.NoError(t, err)
	defer pool.Close()

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
			err := repo.AddSentEmail(ctx, "test@test.com", net.ParseIP("127.0.0.1").To4(), true)
			require.NoError(t, err)

			emails, err := repo.GetSentEmails(ctx)
			require.NoError(t, err)
			require.Len(t, emails, 1)
		},
	)
}
