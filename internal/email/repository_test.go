package email

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"project-a/internal/testutil"
	"testing"
)

func TestEmailRepo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pgContainer, connStr := testutil.SetupTestContainer(ctx, t)

	t.Run(
		"Should add one email row", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

			repo := NewRepo(pool)
			err := repo.AddSentEmail(ctx, "test@test.com", net.ParseIP("127.0.0.1").To4(), true)
			require.NoError(t, err)

			emails, err := repo.GetSentEmails(ctx)
			require.NoError(t, err)
			require.Len(t, emails, 1)
		},
	)
}
