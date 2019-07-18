package transactions

import (
	"errors"
	"time"

	"github.com/rugwirobaker/paypack-backend/app"
)

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("entity already exists")
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid entity format")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
)

//Service defines the transaction service API
type Service interface {
	// RecordTransaction adds a new transaction record
	RecordTransaction(Transaction) (Transaction, error)

	// ViewTransaction retrieves data about the transaction identified with the provided ID,
	ViewTransaction(string) (Transaction, error)

	// ListTransactions retrieves data about subset of transactions that belongs to the
	// user identified by the provided key.
	ListTransactions(uint64, uint64) (TransactionPage, error)

	// ListTransactionsByProperty retrieves data about a subset of transactions that belong to
	// a property identified by the given id. i.e filer by property
	ListTransactionsByProperty(string, uint64, uint64) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that were made using
	// a given method.
	ListTransactionsByMethod(string, uint64, uint64) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that during a given month.
	ListTransactionsByMonth(string, uint64, uint64) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that during a given year.
	ListTransactionsByYear(string, uint64, uint64) (TransactionPage, error)
}

var _ Service = (*transactionsService)(nil)

type transactionsService struct {
	idp   app.IdentityProvider
	store Store
}

//New instantiates a new transaxtions service
func New(idp app.IdentityProvider, store Store) Service {
	return &transactionsService{
		idp:   idp,
		store: store,
	}
}

func (svc *transactionsService) RecordTransaction(transaction Transaction) (Transaction, error) {
	if err := transaction.Validate(); err != nil {
		return Transaction{}, err
	}

	transaction.ID = svc.idp.ID()

	id, err := svc.store.Save(transaction)
	if err != nil {
		return Transaction{}, err
	}
	transaction.ID = id

	transaction.DateRecorded = time.Now()

	return transaction, nil
}

func (svc *transactionsService) ViewTransaction(id string) (Transaction, error) {
	return svc.store.RetrieveByID(id)
}

func (svc *transactionsService) ListTransactions(offset uint64, limit uint64) (TransactionPage, error) {
	return svc.store.RetrieveAll(offset, limit), nil
}

func (svc *transactionsService) ListTransactionsByProperty(id string, offset, limit uint64) (TransactionPage, error) {
	return svc.store.RetrieveByProperty(id, offset, limit), nil
}

func (svc *transactionsService) ListTransactionsByMethod(method string, offset, limit uint64) (TransactionPage, error) {
	return svc.store.RetrieveByMethod(method, offset, limit), nil
}

func (svc *transactionsService) ListTransactionsByMonth(month string, offset, limit uint64) (TransactionPage, error) {
	return TransactionPage{}, nil
}

func (svc *transactionsService) ListTransactionsByYear(year string, offset, limit uint64) (TransactionPage, error) {
	return TransactionPage{}, nil
}
