package mocks

import (
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/models"
	store "github.com/rugwirobaker/paypack-backend/store/transactions"
)

var _ (store.TransactionStore) = (*transactionStoreMock)(nil)

type transactionStoreMock struct {
	mu           sync.Mutex
	counter      uint64
	transactions map[string]models.Transaction
}

// NewTransactionStore creates TransactiobStore mirror
func NewTransactionStore() store.TransactionStore {
	return &transactionStoreMock{
		transactions: make(map[string]models.Transaction),
	}
}

func (str *transactionStoreMock) Save(transaction models.Transaction) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	str.counter++
	transaction.ID = strconv.FormatUint(str.counter, 10)
	str.transactions[transaction.ID] = transaction
	return transaction.ID, nil
}

func (str *transactionStoreMock) RetrieveByID(id string) (models.Transaction, error) {
	str.mu.Lock()
	defer str.mu.Unlock()
	val, ok := str.transactions[id]
	if !ok {
		return models.Transaction{}, models.ErrNotFound
	}

	return val, nil
}

func (str *transactionStoreMock) RetrieveAll(offset, limit uint64) models.TransactionPage {
	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]models.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return models.TransactionPage{}
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

	page := models.TransactionPage{
		Transactions: items,
	}

	return page
}

func (str *transactionStoreMock) RetrieveByProperty(property string, offset, limit uint64) models.TransactionPage {
	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]models.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return models.TransactionPage{}
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.transactions {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Property == property && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := models.TransactionPage{
		Transactions: items,
	}

	return page
}

func (str *transactionStoreMock) RetrieveByMethod(method string, offset, limit uint64) models.TransactionPage {
	str.mu.Lock()
	defer str.mu.Unlock()

	//method = strings.ToLower(method)
	items := make([]models.Transaction, 0)

	if offset < 0 || limit <= 0 {
		return models.TransactionPage{}
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

	page := models.TransactionPage{
		Transactions: items,
	}

	return page
}

func (str *transactionStoreMock) RetrieveByMonth(month string, offset, limit uint64) models.TransactionPage {
	str.mu.Lock()
	defer str.mu.Unlock()

	return models.TransactionPage{}
}

func (str *transactionStoreMock) RetrieveByYear(year string, offset, limit uint64) models.TransactionPage {
	str.mu.Lock()
	defer str.mu.Unlock()

	return models.TransactionPage{}
}
