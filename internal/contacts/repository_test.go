package contacts

import (
	"context"
	"github.com/stretchr/testify/require"
	"project-a/internal/testutil"
	"project-a/internal/user"
	"testing"
)

func TestContactsRepo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pgContainer, connStr := testutil.SetupTestContainer(ctx, t)

	t.Run(
		"Should create contact list", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

			repo := NewRepo(pool)
			userRepo := user.NewUserRepo(pool)
			u, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)
			_, err = repo.CreateContactList(ctx, "My contact list", u.Id)
			require.NoError(t, err)
			list, err := repo.GetContactLists(ctx, u.Id)
			require.NoError(t, err)
			require.Len(t, list, 1)
		},
	)

	t.Run(
		"Should get 1 invitation and invitation match when accepted", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

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
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

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
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

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
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

			repo := NewRepo(pool)
			userRepo := user.NewUserRepo(pool)
			u, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			listName := "My list"
			_, err = repo.CreateContactList(ctx, listName, u.Id)
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

	t.Run(
		"Should create and link contact to list", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

			repo := NewRepo(pool)
			userRepo := user.NewUserRepo(pool)
			inviter, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			invitee, err := userRepo.InsertUser(ctx, "test2", "test2@test.com")
			require.NoError(t, err)

			listName := "My list"
			cl1, err := repo.CreateContactList(ctx, listName, inviter.Id)
			require.NoError(t, err)

			listName2 := "My list"
			cl2, err := repo.CreateContactList(ctx, listName2, invitee.Id)
			require.NoError(t, err)

			contact, err := repo.CreateContact(ctx, inviter.Id, invitee.Id)
			require.NoError(t, err)
			require.Equal(t, contact.User1Id, inviter.Id)
			require.Equal(t, contact.User2Id, invitee.Id)

			err = repo.CreateContactLink(ctx, contact.Id, cl1.Id)
			require.NoError(t, err)

			err = repo.CreateContactLink(ctx, contact.Id, cl2.Id)
			require.NoError(t, err)

			contacts1, err := repo.GetContacts(ctx, inviter.Id)
			require.NoError(t, err)
			require.Len(t, contacts1, 1)
			require.Equal(t, contacts1[0].UserUuid, invitee.Uuid)
			require.Equal(t, contacts1[0].UserId, invitee.Id)
			require.Equal(t, contacts1[0].Username, invitee.Username)

			contacts2, err := repo.GetContacts(ctx, invitee.Id)
			require.NoError(t, err)
			require.Len(t, contacts2, 1)
			require.Equal(t, contacts2[0].UserUuid, inviter.Uuid)
			require.Equal(t, contacts2[0].UserId, inviter.Id)
			require.Equal(t, contacts2[0].Username, inviter.Username)
		},
	)

	t.Run(
		"Delete contact should cascade delete contact link", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)

			repo := NewRepo(pool)
			userRepo := user.NewUserRepo(pool)
			inviter, err := userRepo.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			invitee, err := userRepo.InsertUser(ctx, "test2", "test2@test.com")
			require.NoError(t, err)

			listName := "My list"
			cl1, err := repo.CreateContactList(ctx, listName, inviter.Id)
			require.NoError(t, err)

			listName2 := "My list"
			cl2, err := repo.CreateContactList(ctx, listName2, invitee.Id)
			require.NoError(t, err)

			contact, err := repo.CreateContact(ctx, inviter.Id, invitee.Id)
			require.NoError(t, err)
			require.Equal(t, contact.User1Id, inviter.Id)
			require.Equal(t, contact.User2Id, invitee.Id)

			err = repo.CreateContactLink(ctx, contact.Id, cl1.Id)
			require.NoError(t, err)

			err = repo.CreateContactLink(ctx, contact.Id, cl2.Id)
			require.NoError(t, err)

			contacts1, err := repo.GetContacts(ctx, inviter.Id)
			require.NoError(t, err)
			require.Len(t, contacts1, 1)

			contacts2, err := repo.GetContacts(ctx, invitee.Id)
			require.NoError(t, err)
			require.Len(t, contacts2, 1)

			err = repo.DeleteContact(ctx, contact.Id)
			require.NoError(t, err)

			contacts1, err = repo.GetContacts(ctx, inviter.Id)
			require.NoError(t, err)
			require.Len(t, contacts1, 0)

			contacts2, err = repo.GetContacts(ctx, invitee.Id)
			require.NoError(t, err)
			require.Len(t, contacts2, 0)
		},
	)
}
