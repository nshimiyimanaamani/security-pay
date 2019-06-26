package app

// IdentityProvider defines a persistent identity provider
type IdentityProvider interface {
	// ID generates the unique identifier.
	ID() string
}
