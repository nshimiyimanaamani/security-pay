package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

var _ (transactions.Repository) = (*transactionRepoMock)(nil)

type transactionRepoMock struct {
	mu           sync.Mutex
	counter      uint64
	transactions map[string]transactions.Transaction
}

// NewTransactionStore creates TransactiobStore mirror
func NewTransactionStore() transactions.Repository {
	return &transactionRepoMock{
		transactions: make(map[string]transactions.Transaction),
	}
}

func (str *transactionRepoMock) Save(ctx context.Context, tx transactions.Transaction) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	str.counter++
	tx.ID = strconv.FormatUint(str.counter, 10)
	str.transactions[tx.ID] = tx
	return tx.ID, nil
}

func (str *transactionRepoMock) RetrieveByID(ctx context.Context, id string) (transactions.Transaction, error) {
	str.mu.Lock()
	defer str.mu.Unlock()
	val, ok := str.transactions[id]
	if !ok {
		return transactions.Transaction{}, transactions.ErrNotFound
	}

	return val, nil
}

func (str *transactionRepoMock) RetrieveAll(ctx context.Context, offset, limit uint64) (transactions.TransactionPage, error) {
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

func (str *transactionRepoMock) RetrieveByPeriod(ctx context.Context, offset, limit uint64) (transactions.TransactionPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	return transactions.TransactionPage{}, nil
}
