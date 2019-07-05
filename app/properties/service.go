package properties

import (
	"errors"

	"github.com/rugwirobaker/paypack-backend/app"
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
	ErrNotFound = errors.New("non-existent entity")
)

// Service defines the properties(houses) api
type Service interface {
	// AddProperty adds a unique property entity. Taking a property entity it returns
	// it's unique id an an nil error if the operation is successful
	AddProperty(Property) (Property, error)

	// Update property modifies the mutable fields of the given property entity.
	// it returns the updated property an  nil error if the operation is a success.
	UpdateProperty(Property) error

	// ViewProperty returns a property entity and a nil error if the operation
	// is successful given its unique id.
	ViewProperty(string) (Property, error)

	// ListPropertiesByOwner returns a list of properties that belong to a given owner
	// withing a given range(offset, limit).
	ListPropertiesByOwner(string, uint64, uint64) (PropertyPage, error)

	// ListPropertiesBySector returns a lists of properties in the given sector
	// withing the given range(offset, limit).
	ListPropertiesBySector(string, uint64, uint64) (PropertyPage, error)

	// ListPropertiesByCell returns a lists of properties in the given cell
	// withing the given range(offset, limit).
	ListPropertiesByCell(string, uint64, uint64) (PropertyPage, error)

	// ListPropertiesByVillage returns a lists of properties in the given village
	// withing the given range(offset, limit).
	ListPropertiesByVillage(string, uint64, uint64) (PropertyPage, error)
}

var _ Service = (*propertyService)(nil)

type propertyService struct {
	idp   app.IdentityProvider
	store Store
}

// New instatiates a new property service
func New(idp app.IdentityProvider, store Store) Service {
	return &propertyService{
		idp:   idp,
		store: store,
	}
}

func (svc *propertyService) AddProperty(property Property) (Property, error) {
	if err := property.Validate(); err != nil {
		return Property{}, err
	}

	property.ID = svc.idp.ID()

	id, err := svc.store.Save(property)
	if err != nil {
		return Property{}, err
	}

	property.ID = id
	return property, nil
}

func (svc *propertyService) UpdateProperty(property Property) error {
	if err := property.Validate(); err != nil {
		return err
	}
	return svc.store.UpdateProperty(property)
}

func (svc *propertyService) ViewProperty(id string) (Property, error) {
	return svc.store.RetrieveByID(id)
}

func (svc *propertyService) ListPropertiesByOwner(owner string, offset, limit uint64) (PropertyPage, error) {
	return svc.store.RetrieveByOwner(owner, offset, limit)
}

func (svc *propertyService) ListPropertiesBySector(sector string, offset, limit uint64) (PropertyPage, error) {
	return svc.store.RetrieveBySector(sector, offset, limit)
}

func (svc *propertyService) ListPropertiesByCell(cell string, offset, limit uint64) (PropertyPage, error) {
	return svc.store.RetrieveByCell(cell, offset, limit)
}

func (svc *propertyService) ListPropertiesByVillage(village string, offset, limit uint64) (PropertyPage, error) {
	return svc.store.RetrieveByVillage(village, offset, limit)
}
