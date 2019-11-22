package postgres_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
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

func CleanDB(t *testing.T, tables ...string) {
	t.Helper()
	for _, table := range tables {
		q := fmt.Sprintf("DELETE FROM %s", table)
		db.Exec(q)
	}
}
