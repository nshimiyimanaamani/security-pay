package postgres

import (
	"database/sql"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func seedDB(db *sql.DB) error {
	const op errors.Op = "store/postgres/seedDB"

	var seed = `
		INSERT INTO sectors (sector) VALUES('paypack.test') ON CONFLICT (sector) do nothing;
		INSERT INTO sectors (sector) VALUES('paypack.developers') ON CONFLICT (sector) do nothing;
		INSERT INTO sectors (sector) VALUES('gasabo.kimironko') ON CONFLICT (sector) do nothing;
		INSERT INTO sectors (sector) VALUES('gasabo.remera') ON CONFLICT (sector) do nothing;
	`
	_, err := db.Exec(seed)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
