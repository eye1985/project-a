package contacts

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/interfaces"
	"project-a/internal/models"
)

type contactsRepo struct {
	pool *pgxpool.Pool
}

//go:embed sql/insert_contact_list.sql
var insertContactListSql string

//go:embed sql/insert_contact.sql
var insertContactSql string

//go:embed sql/insert_contact_list_link.sql
var insertContactLinkSql string

//go:embed sql/update_contact_list.sql
var updateContactListSql string

//go:embed sql/delete_contact_list.sql
var deleteContactListSql string

//go:embed sql/get_all_contact_lists.sql
var getAllContactListsSql string

//go:embed sql/get_contact_list.sql
var getContactListSql string

//go:embed sql/get_contact.sql
var getContactSql string

//go:embed sql/insert_invites.sql
var insertInvitesSql string

//go:embed sql/accept_invite.sql
var acceptInviteSql string

//go:embed sql/get_invitations.sql
var getInvitationsSql string

//type Repository interface {
//	GetContactLists(ctx context.Context, userId int64) ([]*List, error)
//	GetContactList(ctx context.Context, contactListId int64) (*List, error)
//	GetContacts(ctx context.Context, userId int64) ([]*Contact, error)
//	CreateContactList(ctx context.Context, name string, userId int64) error
//	CreateContact(ctx context.Context, inviterId int64, inviteeId int64) (*InsertedContact, error)
//	CreateContactLink(ctx context.Context, contactId int64, contactListId int64) error
//	UpdateContactList(ctx context.Context, name string, contactListId int64) error
//	DeleteContactList(ctx context.Context, contactListId int64) error
//	GetInvitations(ctx context.Context, userId int64) ([]*Invitation, error)
//	InviteUser(ctx context.Context, inviterId int64, inviteeId int64) error
//	AcceptInvite(ctx context.Context, uuid uuid.UUID, inviteeId int64) (*AcceptedInvite, error)
//}

func (ulr *contactsRepo) CreateContactList(ctx context.Context, name string, userId int64) error {
	conn, err := ulr.pool.Exec(ctx, insertContactListSql, name, userId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotCreated
	}

	return nil
}

func (ulr *contactsRepo) UpdateContactList(ctx context.Context, name string, contactListId int64) error {
	conn, err := ulr.pool.Exec(ctx, updateContactListSql, name, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotUpdated
	}

	return nil
}

func (ulr *contactsRepo) DeleteContactList(ctx context.Context, contactListId int64) error {
	conn, err := ulr.pool.Exec(ctx, deleteContactListSql, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotDeleted
	}

	return nil
}

func (ulr *contactsRepo) GetContactLists(ctx context.Context, userId int64) ([]*models.List, error) {
	rows, err := ulr.pool.Query(ctx, getAllContactListsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLists []*models.List
	for rows.Next() {
		var userList models.List
		err = rows.Scan(
			&userList.Id,
			&userList.Uuid,
			&userList.Name,
			&userList.CreatedAt,
			&userList.UpdatedAt,
			&userList.UserId,
		)
		if err != nil {
			return nil, err
		}
		userLists = append(userLists, &userList)
	}

	return userLists, nil
}

func (ulr *contactsRepo) GetContactList(ctx context.Context, contactListId int64) (*models.List, error) {
	row := ulr.pool.QueryRow(ctx, getContactListSql, contactListId)
	var userList models.List
	err := row.Scan(
		&userList.Id,
		&userList.Name,
		&userList.CreatedAt,
		&userList.UpdatedAt,
		&userList.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &userList, nil
}

func (ulr *contactsRepo) GetContacts(ctx context.Context, userId int64) ([]*models.Contact, error) {
	rows, err := ulr.pool.Query(ctx, getContactSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := []*models.Contact{}
	for rows.Next() {
		var record models.Contact

		err := rows.Scan(
			&record.UserId,
			&record.UserUuid,
			&record.Username,
			&record.Email,
			&record.ListName,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

func (ulr *contactsRepo) CreateContact(
	ctx context.Context,
	inviterId int64,
	inviteeId int64,
) (*models.InsertedContact, error) {
	row := ulr.pool.QueryRow(ctx, insertContactSql, inviterId, inviteeId)

	var inserted models.InsertedContact
	err := row.Scan(&inserted.Id, &inserted.User1Id, &inserted.User2Id)
	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (ulr *contactsRepo) CreateContactLink(ctx context.Context, contactId int64, contactListId int64) error {
	conn, err := ulr.pool.Exec(ctx, insertContactLinkSql, contactId, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotCreated
	}

	return nil
}

func (ulr *contactsRepo) InviteUser(ctx context.Context, inviterId int64, inviteeId int64) error {
	conn, err := ulr.pool.Exec(ctx, insertInvitesSql, inviterId, inviteeId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return inviteNotCreated
	}

	return nil
}

func (ulr *contactsRepo) AcceptInvite(ctx context.Context, uuid uuid.UUID, inviteeId int64) (
	*models.AcceptedInvite,
	error,
) {
	row := ulr.pool.QueryRow(ctx, acceptInviteSql, uuid, inviteeId)
	var acceptedInvite models.AcceptedInvite

	err := row.Scan(&acceptedInvite.Id, &acceptedInvite.InviterId, &acceptedInvite.InviteeId)
	if err != nil {
		return nil, err
	}

	return &acceptedInvite, nil
}

func (ulr *contactsRepo) GetInvitations(ctx context.Context, userId int64) ([]*models.Invitation, error) {
	rows, err := ulr.pool.Query(ctx, getInvitationsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invitations := []*models.Invitation{}
	for rows.Next() {
		var invitation models.Invitation
		err := rows.Scan(
			&invitation.Id,
			&invitation.Uuid,
			&invitation.InviterId,
			&invitation.InviteeId,
			&invitation.InviterEmail,
			&invitation.InviteeEmail,
			&invitation.Accepted,
		)

		if err != nil {
			return nil, err
		}

		invitations = append(invitations, &invitation)
	}

	return invitations, nil
}

func NewRepo(pool *pgxpool.Pool) interfaces.ContactsRepository {
	return &contactsRepo{
		pool: pool,
	}
}
