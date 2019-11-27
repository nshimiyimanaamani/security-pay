package properties

import (
	"context"
	"errors"

	"github.com/rugwirobaker/paypack-backend/app/identity"
)

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("property already exists")
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid entity format")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent property entity")
)

// nanoid settings
const (
	Alphabet = "1234567890ABCDEF"
	Length   = 8
)

// Service defines the properties(houses) api
type Service interface {
	// RegisterProperty adds a unique property entity. Taking a property entity it returns
	// it's unique id an an nil error if the operation is successful
	RegisterProperty(ctx context.Context, prop Property) (Property, error)

	// Update property modifies the mutable fields of the given property entity.
	// it returns the updated property an  nil error if the operation is a success.
	UpdateProperty(ctx context.Context, prop Property) error

	// RetrieveProperty returns a property entity and a nil error if the operation
	// is successful given its unique id.
	RetrieveProperty(ctx context.Context, uid string) (Property, error)

	// ListPropertiesByOwner returns a list of properties that belong to a given owner
	// withing a given range(offset, limit).
	ListPropertiesByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesBySector returns a lists of properties in the given sector
	// withing the given range(offset, limit).
	ListPropertiesBySector(ctx context.Context, sector string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesByCell returns a lists of properties in the given cell
	// withing the given range(offset, limit).
	ListPropertiesByCell(ctx context.Context, cell string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesByVillage returns a lists of properties in the given village
	// withing the given range(offset, limit).
	ListPropertiesByVillage(ctx context.Context, village string, offset, limit uint64) (PropertyPage, error)
}

var _ Service = (*propertyService)(nil)

type propertyService struct {
	idp  identity.Provider
	repo Repository
}

// New instatiates a new property service
func New(idp identity.Provider, repo Repository) Service {
	return &propertyService{
		idp:  idp,
		repo: repo,
	}
}

func (svc *propertyService) RegisterProperty(ctx context.Context, prop Property) (Property, error) {
	if err := prop.Validate(); err != nil {
		return Property{}, err
	}
	prop.ID = svc.idp.ID()

	id, err := svc.repo.Save(ctx, prop)
	if err != nil {
		return Property{}, err
	}

	prop.ID = id
	return prop, nil
}

func (svc *propertyService) UpdateProperty(ctx context.Context, prop Property) error {
	if err := prop.Validate(); err != nil {
		return err
	}
	return svc.repo.UpdateProperty(ctx, prop)
}

func (svc *propertyService) RetrieveProperty(ctx context.Context, uid string) (Property, error) {
	return svc.repo.RetrieveByID(ctx, uid)
}

func (svc *propertyService) ListPropertiesByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error) {
	return svc.repo.RetrieveByOwner(ctx, owner, offset, limit)
}

func (svc *propertyService) ListPropertiesBySector(ctx context.Context, sector string, offset, limit uint64) (PropertyPage, error) {
	return svc.repo.RetrieveBySector(ctx, sector, offset, limit)
}

func (svc *propertyService) ListPropertiesByCell(ctx context.Context, cell string, offset, limit uint64) (PropertyPage, error) {
	return svc.repo.RetrieveByCell(ctx, cell, offset, limit)
}

func (svc *propertyService) ListPropertiesByVillage(ctx context.Context, village string, offset, limit uint64) (PropertyPage, error) {
	return svc.repo.RetrieveByVillage(ctx, village, offset, limit)
}
