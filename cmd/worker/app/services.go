package app

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/core/auditor"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

func bootAuditor(db *sql.DB) auditor.Service {
	exec := postgres.NewExecutor(db)
	opts := &auditor.Options{Executor: exec}
	return auditor.New(opts)
}
