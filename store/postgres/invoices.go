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

func (repo *invoiceRepository) All(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoices.Retrieve"

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
			earliest_pending_invoices_view 
		WHERE property=$1;
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
			return invoices.Invoice{}, errors.E(op, err, "no invoice found", errors.KindNotFound)
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
