package payment

import "context"

// Client defines the payment initialiion client
type Client interface {
	Pull(context.Context, *TxRequest) (*TxResponse, error)
	Push(context.Context, *TxRequest) (*TxResponse, error)
}
