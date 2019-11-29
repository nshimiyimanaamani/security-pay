package postgres

import (
	"database/sql"

	_ "github.com/lib/pq" // required driver for postgres access
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

//Connect creates and returns a connection to a PostgreSQl instance.
//and returns a non-nil error if there is a failure.
func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}
