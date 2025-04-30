package contacts

import "errors"

var (
	contactsNotCreated = errors.New("contacts not created")
	contactsNotUpdated = errors.New("contacts not updated")
	contactsNotDeleted = errors.New("contacts not deleted")
	contactNotUpdated  = errors.New("contact not updated")
)
