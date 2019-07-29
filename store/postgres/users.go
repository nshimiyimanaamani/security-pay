package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

var _ users.Store = (*userStore)(nil)

type userStore struct {
	db *sql.DB
}

//NewUserStore instanciates a new user Store
func NewUserStore(db *sql.DB) users.Store {
	return &userStore{db}
}

func (str *userStore) Save(user users.User) (string, error) {
	q := `INSERT INTO users (id, email, password, cell) VALUES ($1, $2, $3, $4) RETURNING id`

	if _, err := str.db.Exec(q, user.ID, user.Email, user.Password, user.Cell); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
			return "", users.ErrConflict
		}
		return "", err
	}

	return user.ID, nil
}

func (str *userStore) RetrieveByID(email string) (users.User, error) {
	q := `SELECT password FROM users WHERE email = $1`

	user := users.User{}
	if err := str.db.QueryRow(q, email).Scan(&user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, users.ErrNotFound
		}
		return user, err
	}

	user.Email = email

	return user, nil
}
