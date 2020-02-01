package owners

import "context"

// Repository defines owner storage interface.
type Repository interface {
	// Save adds a new owner to the owner store.
	Save(ctx context.Context, owner Owner) (Owner, error)

	// Updates a given owner entity
	Update(ctx context.Context, owner Owner) error

	// Retrieve retrieves an owner given their id
	Retrieve(ctx context.Context, id string) (Owner, error)

	// Search retrieves an owner given their fname, lname and phone.
	Search(ctx context.Context, owner Owner) (Owner, error)

	// RetrieveAll retrieves a subst of owners.
	RetrieveAll(ctx context.Context, offset, limit uint64) (OwnerPage, error)

	RetrieveByPhone(ctx context.Context, phone string) (Owner, error)
}
