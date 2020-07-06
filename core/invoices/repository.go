package invoices

import "context"

// Repository ...
type Repository interface {
	Generate(ctx context.Context) error
	ListAll(ctx context.Context, property string, months uint) (InvoicePage, error)
	Earliest(ctx context.Context, property string) (Invoice, error)
	ListPending(ctx context.Context, property string, months uint) (InvoicePage, error)
	ListPayed(ctx context.Context, property string, months uint) (InvoicePage, error)
}
