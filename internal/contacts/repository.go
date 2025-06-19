package contacts

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/interfaces"
	"project-a/internal/model"
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

//go:embed sql/delete_contact.sql
var deleteContactSql string

func (cr *contactsRepo) CreateContactList(ctx context.Context, name string, userId int64) (*model.ContactList, error) {
	row := cr.pool.QueryRow(ctx, insertContactListSql, name, userId)
	var contactList model.ContactList
	err := row.Scan(
		&contactList.Id,
		&contactList.Uuid,
		&contactList.Name,
		&contactList.CreatedAt,
		&contactList.UpdatedAt,
		&contactList.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &contactList, nil
}

func (cr *contactsRepo) UpdateContactList(ctx context.Context, name string, contactListId int64) error {
	conn, err := cr.pool.Exec(ctx, updateContactListSql, name, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotUpdated
	}

	return nil
}

func (cr *contactsRepo) DeleteContactList(ctx context.Context, contactListId int64) error {
	conn, err := cr.pool.Exec(ctx, deleteContactListSql, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotDeleted
	}

	return nil
}

func (cr *contactsRepo) GetContactLists(ctx context.Context, userId int64) ([]*model.List, error) {
	rows, err := cr.pool.Query(ctx, getAllContactListsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLists []*model.List
	for rows.Next() {
		var userList model.List
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

func (cr *contactsRepo) GetContactList(ctx context.Context, contactListId int64) (*model.List, error) {
	row := cr.pool.QueryRow(ctx, getContactListSql, contactListId)
	var userList model.List
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

func (cr *contactsRepo) GetContacts(ctx context.Context, userId int64) ([]*model.Contact, error) {
	rows, err := cr.pool.Query(ctx, getContactSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := []*model.Contact{}
	for rows.Next() {
		var record model.Contact

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

func (cr *contactsRepo) CreateContact(
	ctx context.Context,
	inviterId int64,
	inviteeId int64,
) (*model.InsertedContact, error) {
	row := cr.pool.QueryRow(ctx, insertContactSql, inviterId, inviteeId)

	var inserted model.InsertedContact
	err := row.Scan(&inserted.Id, &inserted.User1Id, &inserted.User2Id)
	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (cr *contactsRepo) DeleteContact(ctx context.Context, contactId int64) error {
	conn, err := cr.pool.Exec(ctx, deleteContactSql, contactId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() == 0 {
		return contactNotDeleted
	}

	return nil
}

func (cr *contactsRepo) CreateContactLink(ctx context.Context, contactId int64, contactListId int64) error {
	conn, err := cr.pool.Exec(ctx, insertContactLinkSql, contactId, contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotCreated
	}

	return nil
}

func (cr *contactsRepo) InviteUser(ctx context.Context, inviterId int64, inviteeId int64) error {
	conn, err := cr.pool.Exec(ctx, insertInvitesSql, inviterId, inviteeId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return inviteNotCreated
	}

	return nil
}

func (cr *contactsRepo) AcceptInvite(ctx context.Context, invUUID uuid.UUID, inviteeId int64) (
	*model.AcceptedInvite,
	error,
) {
	row := cr.pool.QueryRow(ctx, acceptInviteSql, invUUID, inviteeId)
	var acceptedInvite model.AcceptedInvite

	err := row.Scan(&acceptedInvite.Id, &acceptedInvite.InviterId, &acceptedInvite.InviteeId)
	if err != nil {
		return nil, err
	}

	return &acceptedInvite, nil
}

func (cr *contactsRepo) GetInvitations(ctx context.Context, userId int64) ([]*model.Invitation, error) {
	rows, err := cr.pool.Query(ctx, getInvitationsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invitations := []*model.Invitation{}
	for rows.Next() {
		var invitation model.Invitation
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
