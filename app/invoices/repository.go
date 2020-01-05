package invoices

import "context"

// Repository ...
type Repository interface {
	Generate(ctx context.Context) error
	ListAll(ctx context.Context, property string, months uint) (InvoicePage, error)
	ListPending(ctx context.Context, property string, months uint) (InvoicePage, error)
	ListPayed(ctx context.Context, property string, months uint) (InvoicePage, error)
}
