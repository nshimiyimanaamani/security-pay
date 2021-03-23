package postgres

import (
	"context"
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/core/scheduler"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (scheduler.Counter) = (*auditableCounter)(nil)

type auditableCounter struct {
	*sql.DB
}

type archivableCounter struct {
	*sql.DB
}

// NewAuditableCounter creates a new postgres backed counter
func NewAuditableCounter(db *sql.DB) scheduler.Counter {
	return &auditableCounter{db}
}

func (c *auditableCounter) Count(ctx context.Context) (int, error) {
	const op errors.Op = "store/postgres/counter.Count"

	var count int

	q := `SELECT count(*) FROM one_month_old_properties_view`

	if err := c.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return 0, errors.E(op, err)
	}
	return count, nil
}

func (c *archivableCounter) Count(ctx context.Context) (int, error) {
	const op errors.Op = "store/postgres/archivable.Count"

	q := `SELECT count(*) FROM one_month_old_properties_view`

	var count int

	if err := c.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return 0, errors.E(op, err)
	}
	return count, nil
}
