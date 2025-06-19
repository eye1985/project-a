package contacts

import "errors"

var (
	contactsNotCreated = errors.New("contacts not created")
	contactsNotUpdated = errors.New("contacts not updated")
	contactsNotDeleted = errors.New("contacts not deleted")
	inviteNotCreated   = errors.New("invite not created")
	contactNotDeleted  = errors.New("contact not deleted")
)
