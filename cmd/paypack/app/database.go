package app

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

// ConnectToDB ...
func ConnectToDB(config *config.DBConfig) (*sql.DB, error) {
	db, err := postgres.Connect(config.URL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
