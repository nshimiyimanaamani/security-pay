package tokens

import "context"

// JWTProvider defines a temporary identifcation provider for users
// this identificatiob is stateless on server side.
type JWTProvider interface {
	// TemporaryKey generates the temporary access token.
	TemporaryKey(ctx context.Context, id string) (string, error)

	// Identity extracts the entity identifier given its secret key.
	Identity(ctx context.Context, token string) (string, error)
}
