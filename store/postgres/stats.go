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
	return &statsRepository{db}
}

func (repo *statsRepository) RetrieveSectorPayRatio(ctx context.Context, sector string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveSectorPayRatio"

	q := `
		select sector, pending, payed from sector_payment_count where sector=$1;
	`
	var label string

	var payed, pending uint64

	err := repo.QueryRowContext(ctx, q, sector).Scan(&label, &pending, &payed)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "sector not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := stats.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}

	return chart, nil
}
func (repo *statsRepository) RetrieveCellPayRatio(ctx context.Context, cell string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveCellPayRatio"

	var label string

	var payed, pending uint64

	q := `
		select cell, pending, payed from cell_payment_count where cell=$1;
	`
	err := repo.QueryRowContext(ctx, q, cell).Scan(&label, &pending, &payed)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "cell not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := stats.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}

	return chart, nil
}

func (repo *statsRepository) RetrieveVillagePayRatio(ctx context.Context, village string) (stats.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.RetrieveVillagePayRatio"

	q := `
		select village, pending, payed from village_payment_count where village=$1;
	`
	var label string

	var payed, pending uint64

	err := repo.QueryRowContext(ctx, q, village).Scan(&label, &pending, &payed)
	if err != nil {
		empty := stats.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "village not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	chart := stats.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}
	return chart, nil
}
