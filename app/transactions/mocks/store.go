package mocks

import (
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

var _ (transactions.Store) = (*transactionStoreMock)(nil)

type transactionStoreMock struct {
	mu           sync.Mutex
	counter      uint64
	transactions map[string]transactions.Transaction
}

// NewTransactionStore creates TransactiobStore mirror
func NewTransactionStore() transactions.Store {
	return &transactionStoreMock{
		transactions: make(map[string]transactions.Transaction),
	}
}

func (str *transactionStoreMock) Save(transaction transactions.Transaction) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	str.counter++
	transaction.ID = strconv.FormatUint(str.counter, 10)
	str.transactions[transaction.ID] = transaction
	return transaction.ID, nil
}

func (str *transactionStoreMock) RetrieveByID(id string) (transactions.Transaction, error) {
	str.mu.Lock()
	defer str.mu.Unlock()
	val, ok := str.transactions[id]
	if !ok {
		return transactions.Transaction{}, transactions.ErrNotFound
	}

	return val, nil
}

func (str *transactionStoreMock) RetrieveAll(offset, limit uint64) (transactions.TransactionPage, error) {
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

func (str *transactionStoreMock) RetrieveByProperty(property string, offset, limit uint64) (transactions.TransactionPage, error) {
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
		if v.MadeFor == property && id >= first && id < last {
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

func (str *transactionStoreMock) RetrieveByMethod(method string, offset, limit uint64) (transactions.TransactionPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	//method = strings.ToLower(method)
	items := make([]transactions.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return transactions.TransactionPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Method == method && id >= first && id < last {
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

func (str *transactionStoreMock) RetrieveByMonth(month string, offset, limit uint64) (transactions.TransactionPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	return transactions.TransactionPage{}, nil
}

func (str *transactionStoreMock) RetrieveByYear(year string, offset, limit uint64) (transactions.TransactionPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	return transactions.TransactionPage{}, nil
}

func (str *transactionStoreMock) UpdateTransaction(tx transactions.Transaction) error {
	return nil
}
