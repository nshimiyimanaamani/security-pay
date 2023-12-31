package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/metrics"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type statsRepository struct {
	*sql.DB
}

// NewStatsRepository ...
func NewStatsRepository(db *sql.DB) metrics.Repository {
	return &statsRepository{db}
}

func (repo *statsRepository) FindSectorRatio(ctx context.Context, sector string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/stats.FindSectorRatio"

	q := `
		select 
			sector, 
			pending_count, 
			payed_count, 
			expired_count
		from 
			sector_payment_metrics 
		where sector=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`
	var label string

	var payed, pending, expired uint64

	err := repo.QueryRowContext(ctx, q, sector, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {
		empty := metrics.Chart{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this sector", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending, "expired": expired},
	}

	return chart, nil
}
func (repo *statsRepository) FindCellRatio(ctx context.Context, cell string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/stats.FindCellRatio"

	var label string

	var payed, pending, expired uint64

	q := `
		select 
			cell, 
			pending_count, 
			payed_count,
			expired_count
		from  
			cell_payment_metrics 
		where cell=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`
	err := repo.QueryRowContext(ctx, q, cell, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {

		empty := metrics.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this cell", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending, "expired": expired},
	}

	return chart, nil
}

func (repo *statsRepository) FindVillageRatio(ctx context.Context, village string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/stats.FindVillageRatio"

	q := `
		select 
			village, 
			pending_count, 
			payed_count,
			expired_count 
		from  
			village_payment_metrics 
		where village=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`
	var label string

	var payed, pending, expired uint64

	err := repo.QueryRowContext(ctx, q, village, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {
		empty := metrics.Chart{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this village", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	chart := metrics.Chart{
		Label: label,
		Data:  map[string]uint64{"payed": payed, "pending": pending, "expired": expired},
	}
	return chart, nil
}

func (repo *statsRepository) ListSectorRatios(ctx context.Context, sector string, y, m uint) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/stats.ListSectorRatios"

	q := `
		select 
			cell, 
			pending_count, 
			payed_count,
			expired_count
		from 
			cell_payment_metrics 
		where sector=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, sector, y, m)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending, expired uint64

		if err := rows.Scan(&label, &pending, &payed, &expired); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data:  map[string]uint64{"payed": payed, "pending": pending, "expired": expired},
		}

		items = append(items, chart)
	}
	return items, nil

}

func (repo *statsRepository) ListCellRatios(ctx context.Context, cell string, y, m uint) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/stats.ListCellRatios"

	q := `
		select 
			village, 
			pending_count, 
			payed_count,
			expired_amount 
		from 
			village_payment_metrics 
		where cell=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, cell, y, m)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending, expired uint64

		if err := rows.Scan(&label, &pending, &payed, &expired); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data:  map[string]uint64{"payed": payed, "pending": pending, "expired": expired},
		}

		items = append(items, chart)
	}
	return items, nil
}

func (repo *statsRepository) FindSectorBalance(ctx context.Context, sector string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindSectorBalance"

	q := `
		select 
			sector, 
			pending_amount, 
			payed_amount,
			expired_amount
		from 
			sector_payment_metrics 
		where sector=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	var label string

	var payed, pending, expired float64

	err := repo.QueryRowContext(ctx, q, sector, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {
		empty := metrics.Chart{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this sector", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data: map[string]uint64{
			"payed":   uint64(payed),
			"pending": uint64(pending),
			"expired": uint64(expired),
		},
	}

	return chart, nil
}

func (repo *statsRepository) FindCellBalance(ctx context.Context, cell string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindCellBalance"

	q := `
		select 
			cell, 
			pending_amount, 
			payed_amount,
			expired_amount
		from 
			cell_payment_metrics 
		where cell=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	var label string

	var payed, pending, expired float64

	err := repo.QueryRowContext(ctx, q, cell, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {
		empty := metrics.Chart{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this cell", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data: map[string]uint64{
			"payed":   uint64(payed),
			"pending": uint64(pending),
			"expired": uint64(expired),
		},
	}

	return chart, nil
}

func (repo *statsRepository) FindVillageBalance(ctx context.Context, cell string, y, m uint) (metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.FindVillageBalance"

	q := `
		select 
			village, 
			pending_amount, 
			payed_amount,
			expired_amount
		from 
			village_payment_metrics 
		where village=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	var label string

	var payed, pending, expired float64

	err := repo.QueryRowContext(ctx, q, cell, y, m).Scan(&label, &pending, &payed, &expired)
	if err != nil {
		empty := metrics.Chart{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "no data found for this village", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	chart := metrics.Chart{
		Label: label,
		Data: map[string]uint64{
			"payed":   uint64(payed),
			"pending": uint64(pending),
			"expired": uint64(expired),
		},
	}

	return chart, nil
}

func (repo *statsRepository) ListSectorBalances(ctx context.Context, sector string, y, m uint) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.ListSectorBalances"

	q := `
		select 
			cell, 
			pending_amount, 
			payed_amount,
			expired_amount
		from 
			cell_payment_metrics 
		where sector=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, sector, y, m)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending, expired float64

		if err := rows.Scan(&label, &pending, &payed, &expired); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data: map[string]uint64{
				"payed":   uint64(payed),
				"pending": uint64(pending),
				"expired": uint64(expired),
			},
		}

		items = append(items, chart)
	}
	return items, nil
}

func (repo *statsRepository) ListCellBalances(ctx context.Context, cell string, y, m uint) ([]metrics.Chart, error) {
	const op errors.Op = "store/postgres/statsRepository.ListCellBalances"

	q := `
		select 
			village, 
			pending_amount, 
			payed_amount, 
			expired_amount
		from 
			village_payment_metrics 
		where cell=$1 and extract(year from period)=$2 and extract(month from period)=$3;
	`

	items := []metrics.Chart{}

	rows, err := repo.QueryContext(ctx, q, cell, y, m)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		var label string

		var payed, pending, expired float64

		if err := rows.Scan(&label, &pending, &payed, &expired); err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}

		chart := metrics.Chart{
			Label: label,
			Data: map[string]uint64{
				"payed":   uint64(payed),
				"pending": uint64(pending),
				"expired": uint64(expired),
			},
		}

		items = append(items, chart)
	}
	return items, nil
}
