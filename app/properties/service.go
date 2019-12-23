package properties

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// var (
// 	// ErrConflict attempt to create an entity with an alreasdy existing id
// 	ErrConflict = errors.New("property already exists")
// 	//ErrInvalidEntity indicates malformed entity specification (e.g.
// 	//invalid username,  password, account).
// 	ErrInvalidEntity = errors.New("invalid property entity")

// 	// ErrPropertyNotFound indicates a non-existent entity request.
// 	ErrPropertyNotFound = errors.New("property not found")

// 	//ErrOwnerNotFound indicates that the referenced owner does not exists yet in the repository
// 	ErrOwnerNotFound = errors.New("owner not found")
// )

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

func (svc *propertyService) RegisterProperty(ctx context.Context, p Property) (Property, error) {
	const op errors.Op = "app/properties/service.RegisterProperty"

	if err := p.Validate(); err != nil {
		return Property{}, errors.E(op, err)
	}

	p.ID = svc.idp.ID()

	property, err := svc.repo.Save(ctx, p)
	if err != nil {
		return Property{}, errors.E(op, err)
	}
	return property, nil
}

func (svc *propertyService) UpdateProperty(ctx context.Context, prop Property) error {
	const op errors.Op = "app/properties/service.UpdateProperty"

	if err := prop.Validate(); err != nil {
		return errors.E(op, err)
	}

	if err := svc.repo.UpdateProperty(ctx, prop); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *propertyService) RetrieveProperty(ctx context.Context, uid string) (Property, error) {
	const op errors.Op = "app/properties/service.RetrieveProperty"

	property, err := svc.repo.RetrieveByID(ctx, uid)
	if err != nil {
		return Property{}, errors.E(op, err)
	}

	return property, nil
}

func (svc *propertyService) ListPropertiesByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListPropertiesByOwner"

	page, err := svc.repo.RetrieveByOwner(ctx, owner, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *propertyService) ListPropertiesBySector(ctx context.Context, sector string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListPropertiesBySector"

	page, err := svc.repo.RetrieveBySector(ctx, sector, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *propertyService) ListPropertiesByCell(ctx context.Context, cell string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListPropertiesByCell"

	page, err := svc.repo.RetrieveByCell(ctx, cell, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *propertyService) ListPropertiesByVillage(ctx context.Context, village string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListPropertiesByVillage"

	page, err := svc.repo.RetrieveByVillage(ctx, village, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}
