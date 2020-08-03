package invoices

import "context"

// Repository ...
type Repository interface {
	// All retrieves all off the invoices of a house
	All(ctx context.Context, property string, months uint) (InvoicePage, error)
	// Earliest retrieves the earliest invoice of house
	Earliest(ctx context.Context, property string) (Invoice, error)
	// Pending retrieves all the pending invoices
	Pending(ctx context.Context, property string, months uint) (InvoicePage, error)
	// Payed retrieves all the payed invoices
	Payed(ctx context.Context, property string, months uint) (InvoicePage, error)
	// Expired retrieves invoices that are due to be archived(have passed payment date)
	Archivable(context.Context) (InvoicePage, error)
}
