package users

import (
	"github.com/rugwirobaker/paypack-backend/models"
)

//UserStore defines the interface to the user datastore
type UserStore interface {
	// RetrieveByID retrieves user by its unique identifier (i.e. email).
	RetrieveByID(string) (models.User, error)

	// Save persists the user account and returns his id. A non-nil error is returned to indicate
	// operation failure.
	Save(models.User) (string, error)
}

var _ UserStore = (*userStore)(nil)

type userStore struct{}

func (str *userStore) RetrieveByID(string) (models.User, error) {
	return models.NewUser("user@gmail.com", "passsword"), nil
}

func (str *userStore) Save(models.User) (string, error) {
	return "user_id", nil
}
