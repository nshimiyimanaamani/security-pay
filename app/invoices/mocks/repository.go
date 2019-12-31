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

func (repo *repository) Retrieve(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.Retrieve"

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
