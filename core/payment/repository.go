package payment

import "context"

// Repository saves validated Transactions to the underlying datastore
type Repository interface {
	//Save a new pending payment to the database
	Save(ctx context.Context, payment Payment) error
	//Find payment by id
	Find(ctx context.Context, id string) (Payment, error)

	//Update the state of an existing payment
	Update(ctx context.Context, payment Payment) error
}
