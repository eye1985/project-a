package email

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"net"
	"project-a/internal/testutil"
	"testing"
)

func TestWithPostgresSQL(t *testing.T) {
	t.Run(
		"Should add one email row", func(t *testing.T) {
			ctx := context.Background()

			pgContainer, pool, err := testutil.SetupTestContainer(ctx)
			testcontainers.CleanupContainer(t, pgContainer)
			require.NoError(t, err)
			defer pool.Close()

			t.Cleanup(
				func() {
					if err := testcontainers.TerminateContainer(pgContainer); err != nil {
						log.Fatalf("failed to terminate container: %s", err)
					}
				},
			)

			repo := NewRepo(pool)
			err = repo.AddSentEmail(ctx, "test@test.com", net.ParseIP("127.0.0.1").To4(), true)
			require.NoError(t, err)

			emails, err := repo.GetSentEmails(ctx)
			require.NoError(t, err)
			require.Len(t, emails, 1)
		},
	)
}
