package app

import (
	"database/sql"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/config"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
)

// PostgresConnect returns a sql.DB connection to postgres
func PostgresConnect(config *config.PostgresConfig) (*sql.DB, error) {
	const op errors.Op = "app.PostgresConnect"

	db, err := postgres.Connect(config.URL)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return db, nil
}
