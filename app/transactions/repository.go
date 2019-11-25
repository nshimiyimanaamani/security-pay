package transactions

import "context"

// Repository defines the api to the transactions data store
type Repository interface {
	// Save adds a new transactiob to the data store returns nil
	// if the operation is successful or otherwise an error.
	Save(ctx context.Context, tx Transaction) (string, error)

	// RetrieveByID retreives a transaction identified by the given id.
	RetrieveByID(ctx context.Context, id string) (Transaction, error)

	// RetrieveAll retrieves the subset of transactions owned by the specified property.
	RetrieveAll(ctx context.Context, offset, limit uint64) (TransactionPage, error)

	// RetrieveByMethod retrieves the subset of transactions that where made using the given method.
	RetrieveByProperty(ctx context.Context, p string, offset, limit uint64) (TransactionPage, error)

	// RetrieveByMonth retrieves the subset of transactions that where made during the given month.
	RetrieveByPeriod(ctx context.Context, offset, limit uint64) (TransactionPage, error)
}
