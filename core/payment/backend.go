package payment

import "context"

// Client defines the payment initialiion client
type Client interface {
	Status(context.Context) (int, error)
	// Auth(appID, appSecret string) (string, error)
	Pull(ctx context.Context, tx Transaction) (Response, error)
	Push(ctx context.Context) error
}
