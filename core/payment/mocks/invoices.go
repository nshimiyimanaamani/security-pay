package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (invoices.Repository) = (*invoicesMock)(nil)

type invoicesMock struct {
	mu       sync.Mutex
	invoices map[string]invoices.Invoice
}

// NewInvoiceRepository creates a mock invoice repository
func NewInvoiceRepository(invs map[string]invoices.Invoice) invoices.Repository {
	return &invoicesMock{
		invoices: invs,
	}
}

func (repo *invoicesMock) Find(ctx context.Context, id uint64) (invoices.Invoice, error) {
	const op errors.Op = "app/invoices/mocks/repository.Find"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	val, ok := repo.invoices[strconv.FormatUint(id, 10)]
	if !ok {
		return invoices.Invoice{}, errors.E(op, "invoice not found", errors.KindNotFound)
	}
	return val, nil
}

func (repo *invoicesMock) All(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.All"

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

func (repo *invoicesMock) Payed(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.Payed"

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

func (repo *invoicesMock) Pending(ctx context.Context, property string, months uint) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.Pending"

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
func (repo *invoicesMock) Unpaid(ctx context.Context, property string) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.Unpaid"

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
			Total: uint(len(repo.invoices)),
		},
	}
	return page, nil
}

func (repo *invoicesMock) Earliest(ctx context.Context, property string) (invoices.Invoice, error) {
	const op errors.Op = "app/invoices/mocks/repository.Earliest"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, vc := range repo.invoices {
		if vc.CreatedAt.Before(time.Now()) {
			return vc, nil
		}
	}
	return invoices.Invoice{}, errors.E(op, "no invoices found", errors.KindNotFound)
}

func (repo *invoicesMock) Archivable(ctx context.Context) (invoices.InvoicePage, error) {
	const op errors.Op = "app/invoices/mocks/repository.Archivable"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return invoices.InvoicePage{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *invoicesMock) Generate(ctx context.Context, id string, amount, months uint) ([]*invoices.Invoice, error) {
	const op errors.Op = "app/invoices/mocks/repository.Generate"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return nil, errors.E(op, "Not implemented", errors.KindNotImplemented)
}
