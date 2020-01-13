package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/stats"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type statsRepository struct {
	*sql.DB
}

// NewStatsRepository ...
func NewStatsRepository(db *sql.DB) stats.Repository {
	return nil
}

func (repo *statsRepository) RetrieveSectorPayRatio(ctx context.Context, sector string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveSectorPayRatio"

	q := `
		SELECT sector, payed, pending FROM sectors_payment_view
	`
	var chart stats.Chart
	err := repo.QueryRowContext(ctx, q, sector).Scan(&chart.Label)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "sector not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return chart, nil
}
func (repo *statsRepository) RetrieveCellPayRatio(ctx context.Context, cell string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveCellPayRatio"

	q := `
		SELECT cell, payed, pending FROM cells_payment_view;
	`
	var chart stats.Chart
	err := repo.QueryRowContext(ctx, q, cell).Scan(&chart.Label)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "cell not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return chart, nil
}

func (repo *statsRepository) RetrieveVillagePayRatio(ctx context.Context, village string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveVillagePayRatio"

	q := `
		SELECT village, payed, pending FROM villages_payment_view
	`
	var chart stats.Chart
	err := repo.QueryRowContext(ctx, q, village).Scan(&chart.Label)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "village not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return chart, nil
}
