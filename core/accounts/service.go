package accounts

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service defines accounts usercases
type Service interface {
	Create(ctx context.Context, acc Account) (Account, error)
	Update(ctx context.Context, acc Account) error
	Retrieve(ctx context.Context, id string) (Account, error)
	List(ctx context.Context, offset, limit uint64) (AccountPage, error)
}

// Options contains accounts.Service creation options
type Options struct {
	Repository Repository
	IDP        identity.Provider
}
type service struct {
	repo Repository
	idp  identity.Provider
}

// New instantiates the accounts.Service
func New(opts *Options) Service {
	return &service{
		repo: opts.Repository,
		idp:  opts.IDP,
	}
}

func (svc *service) Create(ctx context.Context, acc Account) (Account, error) {
	const op errors.Op = "app/accounts/service.Create"

	if err := acc.Validate(); err != nil {
		return Account{}, errors.E(op, err)
	}
	return svc.repo.Save(ctx, acc)
}

func (svc *service) Update(ctx context.Context, acc Account) error {
	const op errors.Op = "app/accounts/service.Update"

	if err := acc.Validate(); err != nil {
		return errors.E(op, err)
	}

	if err := svc.repo.Update(ctx, acc); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) Retrieve(ctx context.Context, id string) (Account, error) {
	const op errors.Op = "app/accounts/service.Retrieve"

	account, err := svc.repo.Retrieve(ctx, id)
	if err != nil {
		return Account{}, errors.E(op, err)
	}
	return account, nil
}

func (svc *service) List(ctx context.Context, offset, limit uint64) (AccountPage, error) {
	const op errors.Op = "app/accounts/service.List"

	page, err := svc.repo.List(ctx, offset, limit)
	if err != nil {
		return AccountPage{}, errors.E(op, err)
	}
	return page, nil
}
