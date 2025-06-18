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
	t.Run(
		"Should create contact list", func(t *testing.T) {
			t.Parallel()
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
			userRepo := user.NewUserRepo(pool)
			u, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)
			err = repo.CreateContactList(ctx, "My contact list", u.Id)
			require.NoError(t, err)
			list, err := repo.GetContactLists(ctx, u.Id)
			require.NoError(t, err)
			require.Len(t, list, 1)
		},
	)

	t.Run(
		"Should get 1 invitation and invitation match when accepted", func(t *testing.T) {
			t.Parallel()
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
			userRepo := user.NewUserRepo(pool)
			inviter, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			invitee, err := userRepo.InsertUser(ctx, "test2", "test2@test.com")
			require.NoError(t, err)

			err = repo.InviteUser(ctx, inviter.Id, invitee.Id)
			require.NoError(t, err)

			inv, err := repo.GetInvitations(ctx, invitee.Id)
			require.NoError(t, err)
			require.Len(t, inv, 1)

			accepted, err := repo.AcceptInvite(ctx, inv[0].Uuid, invitee.Id)
			require.NoError(t, err)

			require.Equal(t, accepted.InviteeId, invitee.Id)
			require.Equal(t, accepted.InviterId, inviter.Id)
		},
	)

	t.Run(
		"Should create one contact", func(t *testing.T) {
			t.Parallel()
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
			userRepo := user.NewUserRepo(pool)
			inviter, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			invitee, err := userRepo.InsertUser(ctx, "test2", "test2@test.com")
			require.NoError(t, err)

			contact, err := repo.CreateContact(ctx, inviter.Id, invitee.Id)
			require.NoError(t, err)
			require.Equal(t, contact.User1Id, inviter.Id)
			require.Equal(t, contact.User2Id, invitee.Id)
		},
	)

	t.Run(
		"Should not allow to have duplicate invites nor bidirectional invites", func(t *testing.T) {
			t.Parallel()
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
			userRepo := user.NewUserRepo(pool)
			inviter, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			invitee, err := userRepo.InsertUser(ctx, "test2", "test2@test.com")
			require.NoError(t, err)

			_, err = repo.CreateContact(ctx, inviter.Id, invitee.Id)
			require.NoError(t, err)
			_, err = repo.CreateContact(ctx, inviter.Id, invitee.Id)
			require.Error(t, err)
			// Bidirectional
			_, err = repo.CreateContact(ctx, invitee.Id, inviter.Id)
			require.Error(t, err)
		},
	)

	t.Run(
		"Should create, update and delete contact lists", func(t *testing.T) {
			t.Parallel()
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
			userRepo := user.NewUserRepo(pool)
			u, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			listName := "My list"
			err = repo.CreateContactList(ctx, listName, u.Id)
			require.NoError(t, err)

			contactLists, err := repo.GetContactLists(ctx, u.Id)
			require.NoError(t, err)
			require.Len(t, contactLists, 1)

			list, err := repo.GetContactList(ctx, contactLists[0].Id)
			require.NoError(t, err)
			require.Equal(t, list.Name, listName)

			updatedListName := "updated"
			err = repo.UpdateContactList(ctx, updatedListName, list.Id)
			require.NoError(t, err)

			list, err = repo.GetContactList(ctx, contactLists[0].Id)
			require.Equal(t, list.Name, updatedListName)

			err = repo.DeleteContactList(ctx, list.Id)
			require.NoError(t, err)

			list2, err := repo.GetContactLists(ctx, u.Id)
			require.NoError(t, err)
			require.Len(t, list2, 0)
		},
	)
}
