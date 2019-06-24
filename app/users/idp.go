package users

// TempIdentityProvider defines a temporary identifcation provider for users
// this identificatiob is stateless on server side.
type TempIdentityProvider interface {
	// TemporaryKey generates the temporary access token.
	TemporaryKey(string) (string, error)

	// Identity extracts the entity identifier given its secret key.
	Identity(string) (string, error)
}