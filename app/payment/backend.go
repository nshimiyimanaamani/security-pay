package payment

import "context"

// Backend defines the payment initialiion client
type Backend interface {
	Status(context.Context) (int, error)
	Auth(appID, appSecret string) (string, error)
	Pull(ctx context.Context, tx Transaction) (Status, error)
}
