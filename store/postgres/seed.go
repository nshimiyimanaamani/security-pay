package postgres

import (
	"database/sql"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func seedDB(db *sql.DB) error {
	const op errors.Op = "store/postgres/seedDB"

	var seed = `
		INSERT INTO sectors (sector) VALUES('paypack.test');
		INSERT INTO sectors (sector) VALUES('paypack.developers');
		INSERT INTO sectors (sector) VALUES('gasabo.kimironko');
		INSERT INTO sectors (sector) VALUES('gasabo.remera');
	`
	_, err := db.Exec(seed)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
