package interfaces

import (
	"context"
	"github.com/google/uuid"
	"project-a/internal/models"
)

type ContactsRepository interface {
	GetContactLists(ctx context.Context, userId int64) ([]*models.List, error)
	GetContactList(ctx context.Context, contactListId int64) (*models.List, error)
	GetContacts(ctx context.Context, userId int64) ([]*models.Contact, error)
	CreateContactList(ctx context.Context, name string, userId int64) error
	CreateContact(ctx context.Context, inviterId int64, inviteeId int64) (*models.InsertedContact, error)
	CreateContactLink(ctx context.Context, contactId int64, contactListId int64) error
	UpdateContactList(ctx context.Context, name string, contactListId int64) error
	DeleteContactList(ctx context.Context, contactListId int64) error
	GetInvitations(ctx context.Context, userId int64) ([]*models.Invitation, error)
	InviteUser(ctx context.Context, inviterId int64, inviteeId int64) error
	AcceptInvite(ctx context.Context, uuid uuid.UUID, inviteeId int64) (*models.AcceptedInvite, error)
}
