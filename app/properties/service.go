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
	// AddProperty adds a unique property entity. Taking a property entity it returns
	// it's unique id an an nil error if the operation is successful
	AddProperty(token string, prop Property) (Property, error)

	// Update property modifies the mutable fields of the given property entity.
	// it returns the updated property an  nil error if the operation is a success.
	UpdateProperty(token string, prop Property) error

	// ViewProperty returns a property entity and a nil error if the operation
	// is successful given its unique id.
	ViewProperty(token, uid string) (Property, error)

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

	// CreateOwner adds a new property adn returns his id if the operation is a success
	CreateOwner(token string, owner Owner) (string, error)

	// Update owner updates the given owner and returns a nil error if
	// the operation is a success.
	UpdateOwner(token string, owner Owner) error

	// ViewOwner returns a owner entity given it's id and returns\
	// a non-nil error the operation failed
	ViewOwner(token, id string) (Owner, error)

	// Listowners returns a subset(offset, limit) of owners and a non-nil error
	ListOwners(token string, offset, limit uint64) (OwnerPage, error)

	// FindOwner owners finds a owner given their fname, lname and phone.
	FindOwner(token, fname, lname, phone string) (Owner, error)
}

var _ Service = (*propertyService)(nil)

type propertyService struct {
	idp        identity.Provider
	owners     OwnerStore
	properties PropertyStore
	auth       AuthBackend
}

// New instatiates a new property service
func New(idp identity.Provider, owners OwnerStore, properties PropertyStore, auth AuthBackend) Service {
	return &propertyService{
		idp:        idp,
		owners:     owners,
		properties: properties,
		auth:       auth,
	}
}

func (svc *propertyService) AddProperty(token string, prop Property) (Property, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Property{}, err
	}
	if err := prop.Validate(); err != nil {
		return Property{}, err
	}
	prop.ID = svc.idp.ID()

	id, err := svc.properties.Save(prop)
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
	return svc.properties.UpdateProperty(prop)
}

func (svc *propertyService) ViewProperty(token, uid string) (Property, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Property{}, err
	}
	return svc.properties.RetrieveByID(uid)
}

func (svc *propertyService) ListPropertiesByOwner(token, owner string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	found, err := svc.owners.Retrieve(owner)
	if err != nil {
		return PropertyPage{}, err
	}
	return svc.properties.RetrieveByOwner(found.ID, offset, limit)
}

func (svc *propertyService) ListPropertiesBySector(token, sector string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.properties.RetrieveBySector(sector, offset, limit)
}

func (svc *propertyService) ListPropertiesByCell(token, cell string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.properties.RetrieveByCell(cell, offset, limit)
}

func (svc *propertyService) ListPropertiesByVillage(token, village string, offset, limit uint64) (PropertyPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return PropertyPage{}, err
	}
	return svc.properties.RetrieveByVillage(village, offset, limit)
}

// CreateOwner adds a new property adn returns his id if the operation is a success
func (svc *propertyService) CreateOwner(token string, owner Owner) (string, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return "", err
	}
	if err := owner.Validate(); err != nil {
		return "", err
	}
	owner.ID = svc.idp.ID()
	return svc.owners.Save(owner)
}

// Update owner updates the given owner and returns a nil error if
// the operation is a success.
func (svc *propertyService) UpdateOwner(token string, owner Owner) error {
	if _, err := svc.auth.Identity(token); err != nil {
		return err
	}
	if err := owner.Validate(); err != nil {
		return err
	}
	return svc.owners.Update(owner)
}

// ViewOwner returns a owner entity given it's id and returns\
// a non-nil error the operation failed
func (svc *propertyService) ViewOwner(token, id string) (Owner, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Owner{}, err
	}
	return svc.owners.Retrieve(id)
}

func (svc *propertyService) ListOwners(token string, offset, limit uint64) (OwnerPage, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return OwnerPage{}, err
	}
	return svc.owners.RetrieveAll(offset, limit)
}

func (svc *propertyService) FindOwner(token, fname, lname, phone string) (Owner, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Owner{}, err
	}
	return svc.owners.FindOwner(fname, lname, phone)
}
