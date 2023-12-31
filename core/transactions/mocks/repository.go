package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (transactions.Repository) = (*repository)(nil)

type repository struct {
	mu           sync.Mutex
	counter      uint64
	transactions map[string]transactions.Transaction
}

// NewRepository creates TransactiobStore mirror
func NewRepository() transactions.Repository {
	return &repository{
		transactions: make(map[string]transactions.Transaction),
	}
}

func (str *repository) Save(ctx context.Context, tx transactions.Transaction) (string, error) {
	const op errors.Op = "app/transactions/mocks/repository.Save"

	str.mu.Lock()
	defer str.mu.Unlock()

	str.counter++
	tx.ID = strconv.FormatUint(str.counter, 10)
	str.transactions[tx.ID] = tx
	return tx.ID, nil
}

func (str *repository) Update(ctx context.Context, tx transactions.Transaction) error {
	const op errors.Op = "app/transactions/mocks/repository.Update"
	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.transactions[tx.ID]; !ok {
		return errors.E(op, "transaction not found", errors.KindNotFound)
	}

	return nil
}

func (str *repository) RetrieveByID(ctx context.Context, id string) (transactions.Transaction, error) {
	const op errors.Op = "app/transactions/mocks/repository.RetrieveByID"

	str.mu.Lock()
	defer str.mu.Unlock()

	empty := transactions.Transaction{}
	val, ok := str.transactions[id]
	if !ok {
		return empty, errors.E(op, "transaction not found", errors.KindNotFound)
	}

	return val, nil
}

func (str *repository) RetrieveAll(ctx context.Context, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "app/transactions/mocks/repository.RetrieveAll"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]transactions.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return transactions.TransactionPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := transactions.TransactionPage{
		Transactions: items,
	}

	return page, nil
}

func (str *repository) RetrieveByProperty(ctx context.Context, p string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "app/transactions/mocks/repository.RetrieveByProperty"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]transactions.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return transactions.TransactionPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.MadeFor == p && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := transactions.TransactionPage{
		Transactions: items,
	}

	return page, nil
}

func (str *repository) RetrieveByPropertyR(ctx context.Context, p string) (transactions.TransactionPage, error) {
	const op errors.Op = "app/transactions/mocks/repository.RetrieveByProperty"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]transactions.Transaction, 0)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		if v.MadeFor == p {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := transactions.TransactionPage{
		Transactions: items,
	}

	return page, nil
}

func (str *repository) RetrieveByMethod(ctx context.Context, m string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "app/transactions/mocks/repository.RetrieveByMethod"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]transactions.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return transactions.TransactionPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Method == m && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := transactions.TransactionPage{
		Transactions: items,
	}

	return page, nil
}
