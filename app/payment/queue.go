package payment

import "context"

// Queue caches unconfirmed  transactions until validation.
type Queue interface {
	//Set adds a new transaction to the cache
	Set(ctx context.Context, tx Transaction) error

	//Unset removes a transaction from the cache for processing
	Unset(ctx context.Context, uid string) (Transaction, error)
}
