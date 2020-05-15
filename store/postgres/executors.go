package postgres

import (
	"context"
	"database/sql"
)

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
	var count int

	q := `SELECT audit_func($1, $2);`

	if err := exec.QueryRowContext(ctx, q, offset, limit).Scan(&count); err != nil {
		return count, err
	}
	return count, nil
}
