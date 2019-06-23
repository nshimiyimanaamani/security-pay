package transactions

import (
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/store/transactions"
)

//Service defines the transaction service API
type Service interface {
	// RecordTransaction adds a new transaction record
	RecordTransaction(models.Transaction) (models.Transaction, error)

	// ViewTransaction retrieves data about the transaction identified with the provided ID,
	ViewTransaction(string) (models.Transaction, error)

	// ListTransactions retrieves data about subset of transactions that belongs to the
	// user identified by the provided key.
	ListTransactions(uint64, uint64) (models.TransactionPage, error)

	// ListTransactionsByProperty retrieves data about a subset of transactions that belong to
	// a property identified by the given id. i.e filer by property
	ListTransactionsByProperty(string, uint64, uint64) (models.TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that were made using
	// a given method.
	ListTransactionsByMethod(string, uint64, uint64) (models.TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that during a given month.
	ListTransactionsByMonth(string, uint64, uint64) (models.TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that during a given year.
	ListTransactionsByYear(string, uint64, uint64) (models.TransactionPage, error)
}

var _ Service = (*transactionsService)(nil)

type transactionsService struct {
	store transactions.Store
}

//New instantiates a new transaxtions service
func New(store transactions.Store) Service {
	return &transactionsService{
		store: store,
	}
}

func (svc *transactionsService) RecordTransaction(transaction models.Transaction) (models.Transaction, error) {
	if err := transaction.Validate(); err != nil {
		return models.Transaction{}, err
	}
	id, err := svc.store.Save(transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	transaction.ID = id
	return transaction, nil
}

func (svc *transactionsService) ViewTransaction(id string) (models.Transaction, error) {
	return svc.store.RetrieveByID(id)
}

func (svc *transactionsService) ListTransactions(offset uint64, limit uint64) (models.TransactionPage, error) {
	return svc.store.RetrieveAll(offset, limit), nil
}

func (svc *transactionsService) ListTransactionsByProperty(id string, offset, limit uint64) (models.TransactionPage, error) {
	return svc.store.RetrieveByProperty(id, offset, limit), nil
}

func (svc *transactionsService) ListTransactionsByMethod(method string, offset, limit uint64) (models.TransactionPage, error) {
	return svc.store.RetrieveByMethod(method, offset, limit), nil
}

/**
 * @todo Implement  ListTransactionsByMonth metnod.
 * @body this method filters transactions by month.
 */
func (svc *transactionsService) ListTransactionsByMonth(month string, offset, limit uint64) (models.TransactionPage, error) {
	return models.TransactionPage{}, nil
}

/**
 * @todo Implement  ListTransactionsByYear metnod.
 * @body this method filters transactions by year.
 */
func (svc *transactionsService) ListTransactionsByYear(year string, offset, limit uint64) (models.TransactionPage, error) {
	return models.TransactionPage{}, nil
}
