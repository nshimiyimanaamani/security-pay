package auth

import "context"

type ctxKey string

const credsKey ctxKey = "auth-creds-context-key"

// SetECredetialsInContext stores Crendetials in the request context
func SetECredetialsInContext(ctx context.Context, c *Credentials) context.Context {
	return context.WithValue(ctx, credsKey, c)
}

// CredentialsFromContext returns Credentials that has been stored in the request context.
// If there is no value for the key or the type assertion fails, it returns a new
// entry from the provided logger
func CredentialsFromContext(ctx context.Context) *Credentials {
	c, ok := ctx.Value(credsKey).(*Credentials)
	if !ok || c == nil {
		return nil
	}
	return c
}
