package users

import (
	"github.com/rugwirobaker/paypack-backend/models"
)

//Store defines the interface to the user datastore
type Store interface {
	// RetrieveByID retrieves user by its unique identifier (i.e. email).
	RetrieveByID(string) (models.User, error)

	// Save persists the user account and returns his id. A non-nil error is returned to indicate
	// operation failure.
	Save(models.User) (string, error)
}
