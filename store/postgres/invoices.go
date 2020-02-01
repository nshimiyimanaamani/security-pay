package postgres

import (
	"context"
	"database/sql"

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

func (repo *invoiceRepository) ListAll(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoiceRepository.Retrieve"

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

func (repo *invoiceRepository) ListPending(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoiceRepository.ListPending"

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

func (repo *invoiceRepository) ListPayed(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "store/postgres/invoiceRepository.ListPayed"

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

func (repo *invoiceRepository) Generate(ctx context.Context) error {
	const op errors.Op = "store/postgres/invoiceRepository.Retrieve"

	return errors.E(op, errors.KindNotImplemented)
}
