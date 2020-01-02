package invoices

import "context"

// Repository ...
type Repository interface {
	Generate(ctx context.Context) error
	Retrieve(ctx context.Context, property string, months uint) (InvoicePage, error)
}
