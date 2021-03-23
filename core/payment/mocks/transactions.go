package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (transactions.Repository) = (*transactionRepoMock)(nil)

type transactionRepoMock struct {
	mu           sync.Mutex
	counter      uint64
	transactions map[string]transactions.Transaction
}

// NewTransactionsRepository creates a memory backed
// transactions.Repository mock.
func NewTransactionsRepository() transactions.Repository {
	return &transactionRepoMock{
		transactions: make(map[string]transactions.Transaction),
	}
}

func (str *transactionRepoMock) Save(ctx context.Context, tx transactions.Transaction) (string, error) {
	const op errors.Op = "core/payment/mocks/repository.Save"

	str.mu.Lock()
	defer str.mu.Unlock()

	str.counter++
	// tx.ID = strconv.FormatUint(str.counter, 10)
	str.transactions[tx.ID] = tx
	return tx.ID, nil
}

func (str *transactionRepoMock) Update(ctx context.Context, tx transactions.Transaction) error {
	const op errors.Op = "core/payment/mocks/repository.Update"

	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.transactions[tx.ID]; !ok {
		return errors.E(op, "transaction not found", errors.KindNotFound)
	}
	return nil
}

func (str *transactionRepoMock) RetrieveByID(ctx context.Context, id string) (transactions.Transaction, error) {
	const op errors.Op = "core/payment/mocks/repository.RetrieveByID("

	str.mu.Lock()
	defer str.mu.Unlock()

	empty := transactions.Transaction{}
	val, ok := str.transactions[id]
	if !ok {
		return empty, errors.E(op, "transaction not found", errors.KindNotFound)
	}

	return val, nil
}

func (str *transactionRepoMock) RetrieveAll(ctx context.Context, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "core/payment/mocks/repository.RetrieveAll"

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

func (str *transactionRepoMock) RetrieveByProperty(ctx context.Context, p string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "core/payment/mocks/repository.RetrieveByProperty"

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

func (str *transactionRepoMock) RetrieveByPropertyR(ctx context.Context, p string) (transactions.TransactionPage, error) {
	const op errors.Op = "core/payment/mocks/repository.RetrieveByPropertyR"

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

func (str *transactionRepoMock) RetrieveByMethod(ctx context.Context, m string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "core/payment/mocks/repository.RetrieveByMethod"

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
