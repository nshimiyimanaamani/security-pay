package app

import (
	"database/sql"

	"github.com/nshimiyimanaamani/paypack-backend/core/archiver"
	"github.com/nshimiyimanaamani/paypack-backend/core/auditor"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
)

// Services ....
type Services struct {
	Auditor  auditor.Service
	Archiver archiver.Service
}

// ProvideServices ...
func ProvideServices(db *sql.DB) *Services {
	return &Services{
		Auditor:  bootAuditor(db),
		Archiver: bootArchiver(db),
	}
}

func bootAuditor(db *sql.DB) auditor.Service {
	exec := postgres.NewExecutor(db)
	opts := &auditor.Options{Executor: exec}
	return auditor.New(opts)
}

func bootArchiver(db *sql.DB) archiver.Service {
	exec := postgres.NewExecutor(db)
	opts := &archiver.Options{Executor: exec}
	return archiver.New(opts)
}
