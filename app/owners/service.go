package owners

import "context"

import "github.com/rugwirobaker/paypack-backend/app/identity"

// Service defines the owners module usecases
type Service interface {
	// Record adds a new property adn returns his id if the operation is a success
	Register(ctx context.Context, owner Owner) (Owner, error)

	// Update owner updates the given owner and returns a nil error if
	// the operation is a success.
	Update(ctx context.Context, owner Owner) error

	// Retrieve returns a owner entity given it's id and returns\
	// a non-nil error the operation failed
	Retrieve(ctx context.Context, id string) (Owner, error)

	// Listowners returns a subset(offset, limit) of owners and a non-nil error
	List(ctx context.Context, offset, limit uint64) (OwnerPage, error)

	// Search owners finds a owner given their fname, lname and phone.
	Search(ctx context.Context, owner Owner) (Owner, error)

	// RetrieveByPhone mobile login
	RetrieveByPhone(ctx context.Context, phone string) (Owner, error)
}

type service struct {
	idp  identity.Provider
	repo Repository
}

// Options ...
type Options struct {
	Idp  identity.Provider
	Repo Repository
}

// New ...
func New(opts *Options) Service {
	return &service{
		idp:  opts.Idp,
		repo: opts.Repo,
	}
}

// Record adds a new property adn returns his id if the operation is a success
func (svc *service) Register(ctx context.Context, owner Owner) (Owner, error) {
	if err := owner.Validate(); err != nil {
		return Owner{}, err
	}
	owner.ID = svc.idp.ID()
	return svc.repo.Save(ctx, owner)
}

// Update owner updates the given owner and returns a nil error if
// the operation is a success.
func (svc *service) Update(ctx context.Context, owner Owner) error {
	if err := owner.Validate(); err != nil {
		return err
	}
	return svc.repo.Update(ctx, owner)
}

// Retrieve returns a owner entity given it's id and returns\
// a non-nil error the operation failed
func (svc *service) Retrieve(ctx context.Context, id string) (Owner, error) {
	return svc.repo.Retrieve(ctx, id)
}

func (svc *service) List(ctx context.Context, offset, limit uint64) (OwnerPage, error) {
	return svc.repo.RetrieveAll(ctx, offset, limit)
}

func (svc *service) Search(ctx context.Context, owner Owner) (Owner, error) {
	return svc.repo.Search(ctx, owner)
}

func (svc *service) RetrieveByPhone(ctx context.Context, phone string) (Owner, error) {
	return svc.repo.RetrieveByPhone(ctx, phone)
}
