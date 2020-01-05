package invoices

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Service ...
type Service interface {
	Create(ctx context.Context) error
	RetrieveAll(ctx context.Context, property string, months uint) (InvoicePage, error)
	RetrievePending(ctx context.Context, property string, months uint) (InvoicePage, error)
	RetrievePayed(ctx context.Context, property string, months uint) (InvoicePage, error)
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

func (svc *service) RetrieveAll(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.RetrieveAll"

	page, err := svc.repo.ListAll(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) RetrievePending(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.RetrievePending"

	page, err := svc.repo.ListPending(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) RetrievePayed(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.RetrievePayed"

	page, err := svc.repo.ListPayed(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) Create(ctx context.Context) error {
	const op errors.Op = "app/invoices/service.Create"

	return errors.E(op, errors.KindNotImplemented)
}
