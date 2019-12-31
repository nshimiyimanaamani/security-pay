package invoices

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Service ...
type Service interface {
	Retrieve(ctx context.Context, property string, months uint) (InvoicePage, error)
}

// Options ...
type Options struct {
	Repo Repository
}

type service struct {
	repo Repository
}

// New ...
func New(opts *Options) Service {
	return &service{
		repo: opts.Repo,
	}
}

func (svc *service) Retrieve(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.Retrieve"

	page, err := svc.repo.Retrieve(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}
