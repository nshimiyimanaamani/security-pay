package auth

import "context"

// Service aggregates Authentication usecases
type Service interface {
	Authenticate(ctx context.Context, token string) error
}
