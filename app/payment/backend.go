package payment

import "context"

// Backend defines the payment initialiion client
type Backend interface {
	Status(context.Context) (int, error)
	Pull(ctx context.Context, tx Transaction) (Status, error)
}
