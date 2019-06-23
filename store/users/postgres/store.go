package postgres

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/store/users"
)

var _ users.Store = (*userStore)(nil)

type userStore struct {
	db *sql.DB
}

//New instanciates a new user Store
func New(db *sql.DB) users.Store {
	return &userStore{db}
}

func (str *userStore) Save(models.User) (string, error) {
	return "user_id", nil
}

func (str *userStore) RetrieveByID(string) (models.User, error) {
	return models.NewUser("user@gmail.com", "passsword"), nil
}
