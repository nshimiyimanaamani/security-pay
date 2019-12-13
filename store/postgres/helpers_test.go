package postgres_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

func saveOwner(t *testing.T, db *sql.DB, owner properties.Owner) (properties.Owner, error) {
	t.Helper()

	q := `INSERT INTO owners (id, fname, lname, phone) VALUES ($1, $2, $3, $4) RETURNING id;`

	_, err := db.Exec(q, &owner.ID, &owner.Fname, &owner.Lname, &owner.Phone)
	if err != nil {
		return properties.Owner{}, err
	}

	return owner, nil
}

func saveTx(t *testing.T, db *sql.DB, tx transactions.Transaction) (transactions.Transaction, error) {
	t.Helper()

	q := `INSERT INTO transactions (id, madefor, madeby, amount,
		 method, date_recorded) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	var empty = transactions.Transaction{}

	_, err := db.Exec(q, tx.ID, tx.MadeFor, tx.MadeBy, tx.Amount, tx.Method, tx.DateRecorded)
	if err != nil {
		return empty, err
	}
	return tx, nil
}

func saveUser(t *testing.T, db *sql.DB, user users.User) (users.User, error) {
	q := `INSERT INTO users (id, username, password, cell, sector, village) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	empty := users.User{}

	if _, err := db.Exec(q, user.ID, user.Username, user.Password, user.Cell, user.Sector, user.Village); err != nil {
		return empty, err
	}
	return user, nil

}

func CleanDB(t *testing.T, tables ...string) {
	t.Helper()
	for _, table := range tables {
		q := fmt.Sprintf("DELETE FROM %s", table)
		db.Exec(q)
	}
}
