package properties

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// nanoid settings
const (
	Alphabet = "1234567890ABCDEF"
	Length   = 8
)

// Service defines the properties(houses) api
type Service interface {
	// Register adds a unique property entity. Taking a property entity it returns
	// it's unique id an an nil error if the operation is successful
	Register(ctx context.Context, prop Property) (Property, error)

	// Update property modifies the mutable fields of the given property entity.
	// it returns the updated property an  nil error if the operation is a success.
	Update(ctx context.Context, prop Property) error

	// Retrieve returns a property entity and a nil error if the operation
	// is successful given its unique id.
	Retrieve(ctx context.Context, uid string) (Property, error)

	// Delete removes a single property given it's id
	// if the operation fails the function returns an error
	Delete(ctx context.Context, uid string) error

	// ListByOwner returns a list of properties that belong to a given owner
	// withing a given range(offset, limit).
	ListByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error)

	// ListByRecorder returns a list of properties as saved by a given user
	//ListByRecorder
	ListByRecorder(ctx context.Context, user string, offset, limit uint64) (PropertyPage, error)

	// ListBySector returns a lists of properties in the given sector
	// withing the given range(offset, limit).
	ListBySector(ctx context.Context, sector string, offset, limit uint64, names string) (PropertyPage, error)

	// ListByCell returns a lists of properties in the given cell
	// withing the given range(offset, limit).
	ListByCell(ctx context.Context, cell string, offset, limit uint64, names string) (PropertyPage, error)

	// ListPropertiesByVillage returns a lists of properties in the given village
	// withing the given range(offset, limit).
	ListByVillage(ctx context.Context, village string, offset, limit uint64, names string) (PropertyPage, error)
}

var _ Service = (*service)(nil)

type service struct {
	idp  identity.Provider
	repo Repository
}

// New instatiates a new property service
func New(idp identity.Provider, repo Repository) Service {
	return &service{
		idp:  idp,
		repo: repo,
	}
}

func (svc *service) Register(ctx context.Context, p Property) (Property, error) {
	const op errors.Op = "app/properties/service.Register"

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

func (svc *service) Update(ctx context.Context, prop Property) error {
	const op errors.Op = "app/properties/service.Update"

	if err := prop.Validate(); err != nil {
		return errors.E(op, err)
	}

	if err := svc.repo.Update(ctx, prop); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) Retrieve(ctx context.Context, uid string) (Property, error) {
	const op errors.Op = "app/properties/service.Retrieve"

	property, err := svc.repo.RetrieveByID(ctx, uid)
	if err != nil {
		return Property{}, errors.E(op, err)
	}

	return property, nil
}

func (svc *service) Delete(ctx context.Context, uid string) error {
	const op errors.Op = "app/properties/service.Delete"

	err := svc.repo.Delete(ctx, uid)
	if err != nil {
		return errors.E(op, err)
	}
	return err
}

func (svc *service) ListByOwner(ctx context.Context, owner string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListByOwner"

	page, err := svc.repo.RetrieveByOwner(ctx, owner, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) ListByRecorder(ctx context.Context, user string, offset, limit uint64) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListByRecorder"

	page, err := svc.repo.RetrieveByRecorder(ctx, user, offset, limit)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) ListBySector(ctx context.Context, sector string, offset, limit uint64, names string) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListBySector"

	page, err := svc.repo.RetrieveBySector(ctx, sector, offset, limit, names)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) ListByCell(ctx context.Context, cell string, offset, limit uint64, names string) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListByCell"

	page, err := svc.repo.RetrieveByCell(ctx, cell, offset, limit, names)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) ListByVillage(ctx context.Context, village string, offset, limit uint64, names string) (PropertyPage, error) {
	const op errors.Op = "app/properties/service.ListPropertiesByVillage"

	page, err := svc.repo.RetrieveByVillage(ctx, village, offset, limit, names)
	if err != nil {
		return PropertyPage{}, errors.E(op, err)
	}
	return page, nil
}
