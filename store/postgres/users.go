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

func (str *userStore) Save(user users.User) (users.User, error) {
	q := `INSERT INTO users (id, username, password, cell, sector, village) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	empty := users.User{}
	if _, err := str.db.Exec(q, user.ID, user.Username, user.Password, user.Cell, user.Sector, user.Village); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
			return empty, users.ErrConflict
		}
		return empty, err
	}

	return user, nil
}

func (str *userStore) RetrieveByID(username string) (users.User, error) {
	q := `SELECT id, password FROM users WHERE username = $1`

	user := users.User{}
	if err := str.db.QueryRow(q, username).Scan(&user.ID, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, users.ErrNotFound
		}
		return user, err
	}

	user.Username = username

	return user, nil
}
