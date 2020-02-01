package identity

// Provider defines a persistent identity provider
type Provider interface {
	// ID generates the unique identifier.
	ID() string
}
