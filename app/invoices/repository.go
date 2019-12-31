package invoices

import "context"

// Repository ...
type Repository interface {
	Retrieve(ctx context.Context, property string, months uint) (InvoicePage, error)
}
