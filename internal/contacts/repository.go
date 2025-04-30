package contacts

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type contactsRepo struct {
	pool *pgxpool.Pool
}

//go:embed sql/insert_contact_list.sql
var insertContactListSql string

//go:embed sql/insert_contact.sql
var insertContactSql string

//go:embed sql/update_contact_list.sql
var updateContactListSql string

//go:embed sql/update_contact.sql
var updateContactSql string

//go:embed sql/delete_contact_list.sql
var deleteContactListSql string

//go:embed sql/get_all_contact_lists.sql
var getAllContactListsSql string

//go:embed sql/get_contact_list.sql
var getContactListSql string

//go:embed sql/get_contact.sql
var getContactSql string

type Repository interface {
	GetContactLists(ctx context.Context, userId int64) ([]*List, error)
	GetContactList(ctx context.Context, contactListId int64) (*List, error)
	GetContactListBySession(ctx context.Context, sessionId string) (*List, error)
	GetContacts(ctx context.Context, userId int64) ([]*Contact, error)
	CreateContactList(ctx context.Context, name string, userId int64) error
	CreateContact(ctx context.Context, userId int64, inviterId int64, contactListId int64, displayName string) (
		*Contact,
		error,
	)
	UpdateContactList(ctx context.Context, name string, contactListId int64) error
	DeleteContactList(ctx context.Context, contactListId int64) error
	UpdateContact(ctx context.Context, hasAccepted bool, contactUuid uuid.UUID, userId int64) error
}

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
	conn, err := ulr.pool.Exec(ctx, updateContactListSql, name, time.Now(), contactListId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactsNotUpdated
	}

	return nil
}

func (ulr *contactsRepo) UpdateContact(
	ctx context.Context,
	hasAccepted bool,
	contactUuid uuid.UUID,
	userId int64,
) error {
	if !hasAccepted {
		return nil
	}

	conn, err := ulr.pool.Exec(ctx, updateContactSql, hasAccepted, time.Now(), contactUuid, userId)
	if err != nil {
		return err
	}

	if conn.RowsAffected() != 1 {
		return contactNotUpdated
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

func (ulr *contactsRepo) GetContactLists(ctx context.Context, userId int64) ([]*List, error) {
	rows, err := ulr.pool.Query(ctx, getAllContactListsSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLists []*List
	for rows.Next() {
		var userList List
		err = rows.Scan(
			&userList.Id,
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

func (ulr *contactsRepo) GetContactList(ctx context.Context, contactListId int64) (*List, error) {
	row := ulr.pool.QueryRow(ctx, getContactListSql, contactListId)
	var userList List
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

func (ulr *contactsRepo) GetContactListBySession(ctx context.Context, sessionId string) (*List, error) {
	row := ulr.pool.QueryRow(ctx, getContactListSql, sessionId)
	var userList List
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

func (ulr *contactsRepo) GetContacts(ctx context.Context, userId int64) ([]*Contact, error) {
	rows, err := ulr.pool.Query(ctx, getContactSql, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var records []*Contact
	for rows.Next() {
		var record Contact

		err := rows.Scan(
			&record.Id,
			&record.Uuid,
			&record.InviterId,
			&record.InviterEmail,
			&record.InviterUsername,
			&record.InviteeId,
			&record.InviteeEmail,
			&record.InviteeUsername,
			&record.HasAccepted,
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
	userId int64,
	inviterId int64,
	contactListId int64,
	displayName string,
) (*Contact, error) {

	var contact Contact
	row := ulr.pool.QueryRow(ctx, insertContactSql, userId, inviterId, displayName, contactListId)
	err := row.Scan(
		&contact.Id,
		&contact.InviteeId,
		&contact.InviterId,
		&contact.HasAccepted,
		&contact.InviteeEmail,
		&contact.InviterEmail,
	)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func NewRepo(pool *pgxpool.Pool) Repository {
	return &contactsRepo{
		pool: pool,
	}
}
