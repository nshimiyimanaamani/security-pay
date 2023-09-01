package postgres

import (
	"context"
	"database/sql"

	"github.com/nshimiyimanaamani/paypack-backend/core/archiver"
	"github.com/nshimiyimanaamani/paypack-backend/core/auditor"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (archiver.Executor) = (*Executor)(nil)
var _ (auditor.Executor) = (*Executor)(nil)

// Executor ...
type Executor struct {
	*sql.DB
}

// NewExecutor ...
func NewExecutor(db *sql.DB) *Executor {
	return &Executor{db}
}

// AuditFunc executes audit query function
func (exec *Executor) AuditFunc(ctx context.Context, offset, limit int) (int, error) {
	const op errors.Op = "store/postgres/Executor.AuditFunc"
	var count int

	q := `SELECT audit_func($1, $2);`

	if err := exec.QueryRowContext(ctx, q, offset, limit).Scan(&count); err != nil {
		return count, errors.E(op, err, errors.KindUnexpected)
	}
	return count, nil
}

// ArchiveFunc executes audit query function
func (exec *Executor) ArchiveFunc(ctx context.Context, offset, limit int) (int, error) {
	const op errors.Op = "store/postgres/Executor.AuditFunc"

	var count int

	q := `SELECT archive_func();`

	if err := exec.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return count, errors.E(op, err, errors.KindUnexpected)
	}
	return count, nil
}
