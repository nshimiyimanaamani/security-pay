package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type invoiceRepository struct {
	*sql.DB
}

// NewInvoiceRepository ...
func NewInvoiceRepository(db *sql.DB) invoices.Repository {
	return &invoiceRepository{db}
}

func (repo *invoiceRepository) Find(ctx context.Context, id uint64) (invoices.Invoice, error) {
	const op errors.Op = "store/postgres/invoices.Find"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM invoices WHERE id=$1 
	`
	var invoice = invoices.Invoice{}

	tx, err := repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	if err != nil {
		return invoices.Invoice{}, errors.E(op, err, errors.KindUnexpected)
	}

	err = tx.QueryRow(q, id).Scan(
		&invoice.ID,
		&invoice.Amount,
		&invoice.Property,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		var empty = invoices.Invoice{}

		if err := tx.Rollback(); err != nil {
			return empty, errors.E(op, err, errors.KindUnexpected)
		}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "invoice not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	if err := tx.Commit(); err != nil {
		return invoices.Invoice{}, errors.E(op, errors.E(op, errors.KindUnexpected))
	}
	return invoice, nil
}

func (repo *invoiceRepository) All(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.All"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM 
			invoices 
		WHERE 
			property=$1 
		AND 
			created_at >= DATE_TRUNC('month', CURRENT_TIMESTAMP) - INTERVAL '1 month' * $2
		ORDER BY created_at DESC
	`

	items := []invoices.Invoice{}

	rows, err := repo.Query(q, property, months)
	if err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := invoices.Invoice{}

		if err := rows.Scan(&c.ID, &c.Amount, &c.Property, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM invoices WHERE property=$1;`

	var total uint

	if err := repo.QueryRow(q, property).Scan(&total); err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  total,
			Months: months,
		},
	}
	return page, nil
}

func (repo *invoiceRepository) Pending(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.Pending"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM 
			invoices 
		WHERE 
			property=$1 
		AND
			status='pending' 
		AND 
			created_at >= DATE_TRUNC('month', CURRENT_TIMESTAMP) - INTERVAL '1 month' * $2
		ORDER BY created_at DESC
	`

	items := []invoices.Invoice{}

	rows, err := repo.Query(q, property, months)
	if err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := invoices.Invoice{}

		if err := rows.Scan(&c.ID, &c.Amount, &c.Property, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM invoices WHERE property=$1 AND status='pending';`

	var total uint

	if err := repo.QueryRow(q, property).Scan(&total); err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  total,
			Months: months,
		},
	}
	return page, nil
}

func (repo *invoiceRepository) Payed(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.Payed"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM 
			invoices 
		WHERE 
			property=$1
		AND 
			status='payed' 
		AND 
			created_at >= DATE_TRUNC('month', CURRENT_TIMESTAMP) - INTERVAL '1 month' * $2
		ORDER BY created_at DESC
	`

	items := []invoices.Invoice{}

	rows, err := repo.Query(q, property, months)
	if err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := invoices.Invoice{}

		if err := rows.Scan(&c.ID, &c.Amount, &c.Property, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM invoices WHERE property=$1 AND status='payed';`

	var total uint

	if err := repo.QueryRow(q, property).Scan(&total); err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  total,
			Months: months,
		},
	}
	return page, nil
}

func (repo *invoiceRepository) Earliest(ctx context.Context, property string) (invoices.Invoice, error) {
	const op errors.Op = "store/postgres/invoiceRepository.Earliest"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM 
			invoices 
		WHERE
			property=$1 AND status='pending'
		AND  
			DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE) ;
	`
	var invoice invoices.Invoice

	err := repo.QueryRow(q, property).Scan(
		&invoice.ID,
		&invoice.Amount,
		&invoice.Property,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return invoices.Invoice{}, errors.E(op, "ntawenda ufite wo kwishyura uku kwezi", errors.KindNotFound)
		}
		return invoices.Invoice{}, errors.E(op, err, errors.KindUnexpected)
	}
	return invoice, nil
}

func (repo *invoiceRepository) Archivable(ctx context.Context) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.Archivable"

	q := `
		select 
			id,
			amount,
			property,
			status,
			created_at,
			updated_at
		FROM 
			invoices
		WHERE 
			status='pending' 
		AND 
			created_at < date_trunc('month', now())::date
		ORDER BY id
	`

	var items []invoices.Invoice

	rows, err := repo.Query(q)
	if err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := invoices.Invoice{}

		if err := rows.Scan(&c.ID, &c.Amount, &c.Property, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM invoices WHERE status='pending' AND created_at < date_trunc('month', now())::date;`

	var total uint

	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total: total,
		},
	}
	return page, nil
}

// Unpaid returns all unpaid invoices for a given property
func (repo *invoiceRepository) Unpaid(ctx context.Context, property string) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.Unpaid"

	q := `
		SELECT 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at 
		FROM 
			invoices 
		WHERE 
			property=$1 
		AND
			status='pending' 
		AND 
			created_at < DATE_TRUNC('month', CURRENT_DATE)
		ORDER BY created_at DESC
	`

	items := []invoices.Invoice{}

	rows, err := repo.Query(q, property)
	if err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := invoices.Invoice{}

		if err := rows.Scan(&c.ID, &c.Amount, &c.Property, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT 
			COUNT(*),
			COALESCE(sum(amount), 0.0) as total_unpaid
		FROM 
			invoices 
		WHERE 
			property=$1 AND status='pending' AND created_at < DATE_TRUNC('month', CURRENT_DATE);`

	var (
		total       uint
		totalUnpaid float64
	)

	if err := repo.QueryRow(q, property).Scan(&total, &totalUnpaid); err != nil {
		return invoices.InvoicePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:       total,
			TotalAmount: totalUnpaid,
		},
	}
	return page, nil
}

// Generate generates a new invoice for a given months count from the current month on wards with the given amount and property id
func (repo *invoiceRepository) Generate(ctx context.Context, property string, amount, months uint) ([]*invoices.Invoice, error) {
	const op errors.Op = "store/postgres/invoices.Generate"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	defer tx.Rollback()

	out := make([]*invoices.Invoice, 0)

	selectQuery := `
		SELECT
			id, amount, property, status, created_at, updated_at
		FROM
			invoices
		WHERE
			property=$1
		AND
			DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)
	`

	invoice := new(invoices.Invoice)
	if err := tx.QueryRow(
		selectQuery,
		property,
	).Scan(
		&invoice.ID,
		&invoice.Amount,
		&invoice.Property,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	); err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	if invoice.Status == invoices.Pending {
		out = append(out, invoice)
	}

	insertQuery := `
		INSERT INTO invoices 
			(amount, property, status, created_at, updated_at)
		SELECT
			$1::numeric, 
			$2::text,
			'pending', 
			DATE_TRUNC('month', CURRENT_DATE) + interval '1 month' * s.a, 
			DATE_TRUNC('month', CURRENT_DATE) + interval '1 month' * s.a
		FROM 
			generate_series(1, $3::int) s(a)
		ON CONFLICT DO NOTHING RETURNING 
			id, 
			amount, 
			property, 
			status, 
			created_at, 
			updated_at
	`

	rows, err := tx.Query(insertQuery, amount/months, property, months)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && errInvalid == pqErr.Code.Name() {
			return nil, errors.E(op, err, errors.KindNotFound)
		}
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		invoice := new(invoices.Invoice)
		err := rows.Scan(
			&invoice.ID,
			&invoice.Amount,
			&invoice.Property,
			&invoice.Status,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, errors.E(op, err, errors.KindUnexpected)
		}
		out = append(out, invoice)
	}

	return out, tx.Commit()
}
