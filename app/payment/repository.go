package payment

import "context"

// Repository saves validated Transactions to the underlying datastore
type Repository interface {
	Save(ctx context.Context, tx Transaction) error
	RetrieveCode(ctx context.Context, code string) (string, error)
}
