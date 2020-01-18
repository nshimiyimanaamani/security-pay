package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/metrics"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type statsRepository struct {
	*sql.DB
}

// NewStatsRepository ...
func NewStatsRepository(db *sql.DB) metrics.Repository {
	return &statsRepository{db}
}

func (repo *statsRepository) FindSectorRatio(ctx context.Context, sector string) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindSectorRatio"

	q := `
		select sector, pending, payed from sector_payment_ratio where sector=$1;
	`
	var label string

	var payed, pending uint64

	err := repo.QueryRowContext(ctx, q, sector).Scan(&label, &pending, &payed)
	if err != nil {
		empty := metrics.Chart{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "sector not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}

	return chart, nil
}
func (repo *statsRepository) FindCellRatio(ctx context.Context, cell string) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindCellRatio"

	var label string

	var payed, pending uint64

	q := `
		select cell, pending, payed from  cell_payment_ratio where cell=$1;
	`
	err := repo.QueryRowContext(ctx, q, cell).Scan(&label, &pending, &payed)
	if err != nil {

		empty := metrics.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "cell not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}

	return chart, nil
}

func (repo *statsRepository) FindVillageRatio(ctx context.Context, village string) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindVillageRatio"

	q := `
		select village, pending, payed from  village_payment_ratio where village=$1;
	`
	var label string

	var payed, pending uint64

	err := repo.QueryRowContext(ctx, q, village).Scan(&label, &pending, &payed)
	if err != nil {
		empty := metrics.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "village not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending},
	}
	return chart, nil
}

func (repo *statsRepository) ListSectorRatios(ctx context.Context, sector string) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.ListSectorRatios"

	q := `
		select cell, pending, payed from cell_payment_ratio where sector=$1;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, sector)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending uint64

		if err := rows.Scan(&label, &pending, &payed); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data:  map[string]uint64{"payed": payed, "pending": pending},
		}

		items = append(items, chart)
	}
	return items, nil

}

func (repo *statsRepository) ListCellRatios(ctx context.Context, cell string) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.ListCellRatios"

	q := `
		select village, pending, payed from village_payment_ratio where cell=$1;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, cell)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending uint64

		if err := rows.Scan(&label, &payed, &pending); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data:  map[string]uint64{"payed": payed, "pending": pending},
		}

		items = append(items, chart)
	}
	return items, nil
}
