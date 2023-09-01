package invoices

import (
	"context"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// Service ...
type Service interface {
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

	page, err := svc.repo.All(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) RetrievePending(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.RetrievePending"

	page, err := svc.repo.Pending(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) RetrievePayed(ctx context.Context, property string, months uint) (InvoicePage, error) {
	const op errors.Op = "app/invoices/service.RetrievePayed"

	page, err := svc.repo.Payed(ctx, property, months)
	if err != nil {
		return InvoicePage{}, errors.E(op, err)
	}
	return page, nil
}
