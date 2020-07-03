package payment

import "context"

// Client defines the payment initialiion client
type Client interface {
	Pull(ctx context.Context, tx Transaction) (Response, error)
	Push(ctx context.Context, tx Transaction) (Response, error)
}
