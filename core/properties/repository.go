package properties

import "context"

// Repository defines the api to the properties data store
type Repository interface {
	// Save adds a new transactiob to the data store returns nil
	// if the operation is successful or otherwise an error.
	Save(ctx context.Context, p Property) (Property, error)

	// Update the given property entity's mutable fields.
	Update(ctx context.Context, p Property) error

	//Delete removes a single entity fron the underlying store.
	Delete(ctx context.Context, uid string) error

	// RetrieveByID retrieves a property entity  given it's unique id.
	RetrieveByID(ctx context.Context, uid string) (Property, error)

	// RetrieveByOwner retrieves the subset of properties that where made using the given owner.
	RetrieveByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error)

	// RetrieveByOwner retrieves the subset of properties that where made using the given a user.
	RetrieveByRecorder(ctx context.Context, user string, offset, limit uint64) (PropertyPage, error)

	// RetrieveBySector retrieves the subset of properties within a given Sector.
	RetrieveBySector(ctx context.Context, sector string, offset, limit uint64) (PropertyPage, error)

	// RetrieveByCell retrieves the subset of properteis within a given Cell.
	RetrieveByCell(ctx context.Context, cell string, offset, limit uint64) (PropertyPage, error)

	// RetrieveByVillage retrieves the subset of properties within a given Village.
	RetrieveByVillage(ctx context.Context, Village string, offset, limit uint64) (PropertyPage, error)

	// Auditable counts the number of properties that need an invoice
	//Auditable(ctx context.Context) (int, error)

	//RetrieveBYNames
	RetrieveByNames(ctx context.Context, names string, offset, limit uint64) (PropertyPage, error)
}
