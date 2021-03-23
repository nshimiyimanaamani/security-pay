package auth

import "context"

// Repository is the interface to logins database
type Repository interface {
	Retrieve(ctx context.Context, username string) (Credentials, error)
}
