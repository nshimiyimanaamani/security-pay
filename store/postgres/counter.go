package postgres

import (
	"context"
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/core/scheduler"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (scheduler.Auditable) = (*counter)(nil)

type counter struct {
	*sql.DB
}

// NewCounter creates a new postgres backed counter
func NewCounter(db *sql.DB) scheduler.Auditable {
	return &counter{db}
}

func (c *counter) Count(ctx context.Context) (int, error) {
	const op errors.Op = "store/postgres/counter.Count"

	var count int

	q := `SELECT count(*) FROM one_month_old_properties_view`

	if err := c.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return 0, errors.E(op, err)
	}
	return count, nil
}
