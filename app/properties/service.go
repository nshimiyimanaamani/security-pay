package properties

import (
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
	ErrNotFound = errors.New("non-existent entity")
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
	RegisterProperty(token string, prop Property) (Property, error)

	// Update property modifies the mutable fields of the given property entity.
	// it returns the updated property an  nil error if the operation is a success.
	UpdateProperty(token string, prop Property) error

	// RetrieveProperty returns a property entity and a nil error if the operation
	// is successful given its unique id.
	RetrieveProperty(token, uid string) (Property, error)

	// ListPropertiesByOwner returns a list of properties that belong to a given owner
	// withing a given range(offset, limit).
	ListPropertiesByOwner(token, owner string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesBySector returns a lists of properties in the given sector
	// withing the given range(offset, limit).
	ListPropertiesBySector(token, sector string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesByCell returns a lists of properties in the given cell
	// withing the given range(offset, limit).
	ListPropertiesByCell(token, cell string, offset, limit uint64) (PropertyPage, error)

	// ListPropertiesByVillage returns a lists of properties in the given village
	// withing the given range(offset, limit).
	ListPropertiesByVillage(token, village string, offset, limit uint64) (PropertyPage, error)
}

var _ Service = (*propertyService)(nil)

type propertyService struct {
	idp  identity.Provider
	repo Repository
	auth AuthBackend
}

// New instatiates a new property service
func New(idp identity.Provider, repo Repository, auth AuthBackend) Service {
	return &propertyService{
		idp:  idp,
		repo: repo,
		auth: auth,
	}
}

func (svc *propertyService) RegisterProperty(token string, prop Property) (Property, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Property{}, err
	}
	if err := prop.Validate(); err != nil {
		return Property{}, err
	}
	prop.ID = svc.idp.ID()

	id, err := svc.repo.Save(prop)
	if err != nil {
		return Property{}, err
	}

	prop.ID = id
	return prop, nil
}

func (svc *propertyService) UpdateProperty(token string, prop Property) error {
	if _, err := svc.auth.Identity(token); err != nil {
		return err
	}
	if err := prop.Validate(); err != nil {
		return err
	}
	return svc.repo.UpdateProperty(prop)
}

func (svc *propertyService) RetrieveProperty(token, uid string) (Property, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Property{}, err
	}
	return svc.repo.RetrieveByID(uid)
}

func (svc *propertyService) ListPropertiesByOwner(token, owner string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.repo.RetrieveByOwner(owner, offset, limit)
}

func (svc *propertyService) ListPropertiesBySector(token, sector string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.repo.RetrieveBySector(sector, offset, limit)
}

func (svc *propertyService) ListPropertiesByCell(token, cell string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.repo.RetrieveByCell(cell, offset, limit)
}

func (svc *propertyService) ListPropertiesByVillage(token, village string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.repo.RetrieveByVillage(village, offset, limit)
}
