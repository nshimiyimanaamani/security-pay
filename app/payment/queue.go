package payment

import "context"

// Queue caches unconfirmed  transactions until validation.
type Queue interface {
	//Set adds a new transaction to the cache
	Set(ctx context.Context, tx Transaction) error

	//Pop removes a transaction from the cache for processing
	Get(ctx context.Context, uid string) (Transaction, error)

	//Remove an already processed transaction from the queue
	Remove(ctx context.Context, uid string) error
}
