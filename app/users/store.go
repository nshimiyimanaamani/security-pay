package users

//Store defines the interface to the user datastore
type Store interface {
	// RetrieveByID retrieves user by its unique identifier (i.e. email).
	RetrieveByID(string) (User, error)

	// Save persists the user account and returns his id. A non-nil error is returned to indicate
	// operation failure.
	Save(User) (string, error)
}
