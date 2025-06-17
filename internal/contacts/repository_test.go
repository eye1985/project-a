package contacts

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"project-a/internal/testutil"
	"project-a/internal/user"
	"testing"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, pool, err := testutil.SetupTestContainer(ctx, "contacts-snapshot")
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)
	defer pool.Close()

	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	t.Run("Should create contact list", func (t *testing.T){
		t.Cleanup(
			func() {
				err := pgContainer.Restore(ctx)
				require.NoError(t, err)
			},
		)

		repo := NewRepo(pool)
		userRepo := user.NewUserRepo(pool)
		u, err := userRepo.InsertUser(ctx, "test", "test@test.com")
		require.NoError(t, err)
		err = repo.CreateContactList(ctx, "My contact list", u.Id)
		require.NoError(t, err)
		list, err := repo.GetContactLists(ctx, u.Id)
		require.NoError(t, err)
		require.Len(t, list, 1)
	})
}
