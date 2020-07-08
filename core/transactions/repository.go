package transactions

import "context"

// Repository defines the api to the transactions data store
type Repository interface {
	// Save adds a new transactiob to the data store returns nil
	// if the operation is successful or otherwise an error.
	Save(ctx context.Context, tx Transaction) (string, error)

	// Update changes the state of an existing transaction
	// if the operation is successful or otherwise an error.
	Update(ctx context.Context, tx Transaction) error

	// RetrieveByID retreives a transaction identified by the given id.
	RetrieveByID(ctx context.Context, id string) (Transaction, error)

	// RetrieveAll retrieves the subset of transactions owned by the specified property.
	RetrieveAll(ctx context.Context, offset, limit uint64) (TransactionPage, error)

	// RetrieveByMethod retrieves the subset of transactions that where made using the given method.
	RetrieveByProperty(ctx context.Context, p string, offset, limit uint64) (TransactionPage, error)

	// RetrieveByMethod retrieves the subset of transactions that where made using the given method.
	//(rigobert's wishes are orders) no need to test here as long as it is delivered
	//You can retrieve and send 40000 objects to the client(Mobile) it's where thought.
	RetrieveByPropertyR(ctx context.Context, p string) (TransactionPage, error)

	// RetrieveByMethodretrieves the subset of transactions that where made during the given month.
	RetrieveByMethod(ctx context.Context, m string, offset, limit uint64) (TransactionPage, error)
}
