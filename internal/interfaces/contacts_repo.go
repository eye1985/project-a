package interfaces

import (
	"context"
	"github.com/google/uuid"
	"project-a/internal/model"
)

type ContactsRepository interface {
	GetContactLists(ctx context.Context, userId int64) ([]*model.List, error)
	GetContactList(ctx context.Context, contactListId int64) (*model.List, error)
	GetContacts(ctx context.Context, userId int64) ([]*model.Contact, error)
	CreateContactList(ctx context.Context, name string, userId int64) (*model.ContactList, error)
	CreateContact(ctx context.Context, inviterId int64, inviteeId int64) (*model.InsertedContact, error)
	CreateContactLink(ctx context.Context, contactId int64, contactListId int64) error
	UpdateContactList(ctx context.Context, name string, contactListId int64) error
	DeleteContactList(ctx context.Context, contactListId int64) error
	DeleteContact(ctx context.Context, contactId int64) error
	GetInvitations(ctx context.Context, userId int64) ([]*model.Invitation, error)
	InviteUser(ctx context.Context, inviterId int64, inviteeId int64) error
	AcceptInvite(ctx context.Context, invUUID uuid.UUID, inviteeId int64) (*model.AcceptedInvite, error)
}
