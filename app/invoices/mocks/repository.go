package mocks

import (
	"context"
	"sort"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (invoices.Repository) = (*repository)(nil)

type repository struct {
	mu       sync.Mutex
	invoices map[string]invoices.Invoice
}

// NewRepository ...
func NewRepository(invs map[string]invoices.Invoice) invoices.Repository {
	return &repository{
		invoices: invs,
	}
}

func (repo *repository) ListAll(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.ListAll"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]invoices.Invoice, 0)

	val, ok := repo.invoices[property]
	if !ok {
		return invoices.InvoicePage{}, errors.E(op, "property doesn't exists", errors.KindNotFound)
	}

	items = append(items, val)

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  uint(len(repo.invoices)),
			Months: months,
		},
	}
	return page, nil
}

func (repo *repository) ListPayed(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.ListPayed"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]invoices.Invoice, 0)

	val, ok := repo.invoices[property]
	if !ok {
		return invoices.InvoicePage{}, errors.E(op, "property doesn't exists", errors.KindNotFound)
	}

	items = append(items, val)

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  uint(len(repo.invoices)),
			Months: months,
		},
	}
	return page, nil
}

func (repo *repository) ListPending(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.ListPending"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]invoices.Invoice, 0)

	val, ok := repo.invoices[property]
	if !ok {
		return invoices.InvoicePage{}, errors.E(op, "property doesn't exists", errors.KindNotFound)
	}

	items = append(items, val)

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := invoices.InvoicePage{
		Invoices: items,
		PageMetadata: invoices.PageMetadata{
			Total:  uint(len(repo.invoices)),
			Months: months,
		},
	}
	return page, nil
}

func (repo *repository) Generate(ctx context.Context) error {
	const op errors.Op = "app/invoices/mocks/repository.Retrieve"

	return errors.E(op, errors.KindNotImplemented)
}
