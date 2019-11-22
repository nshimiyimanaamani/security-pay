package transactions

import (
	"errors"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/identity"
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

// nanoid conf
var (
	Alphabet = "1234567890abcdefghijklmnopqrstuvwxyz"
	Length   = 20
)

//Service defines the transaction service API
type Service interface {
	// RecordTransaction adds a new transaction record
	RecordTransaction(token string, trx Transaction) (Transaction, error)

	// ViewTransaction retrieves data about the transaction identified with the provided ID,
	ViewTransaction(uid string) (Transaction, error)

	// ListTransactions retrieves data about subset of transactions that belongs to the
	// user identified by the provided key.
	ListTransactions(offset, limit uint64) (TransactionPage, error)

	// ListTransactionsByProperty retrieves data about a subset of transactions that belong to
	// a property identified by the given id. i.e filer by property
	ListTransactionsByProperty(prop string, offset, limit uint64) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that were made using
	// a given method.
	ListTransactionsByMethod(meth string, offset, limit uint64) (TransactionPage, error)

	// // ListTransactionByDate retrieves data about a subset of transactions that during a given month.
	// ListTransactionsByMonth(string, uint64, uint64) (TransactionPage, error)

	// // ListTransactionByDate retrieves data about a subset of transactions that during a given year.
	// ListTransactionsByYear(string, uint64, uint64) (TransactionPage, error)
}

var _ Service = (*transactionsService)(nil)

type transactionsService struct {
	idp   identity.Provider
	store Store
	auth  AuthBackend
}

//New instantiates a new transaxtions service
func New(idp identity.Provider, store Store, auth AuthBackend) Service {
	return &transactionsService{
		idp:   idp,
		store: store,
		auth:  auth,
	}
}

func (svc *transactionsService) RecordTransaction(token string, trx Transaction) (Transaction, error) {
	if _, err := svc.auth.Identity(token); err != nil {
		return Transaction{}, err
	}
	if err := trx.Validate(); err != nil {
		return Transaction{}, err
	}

	trx.ID = svc.idp.ID()

	id, err := svc.store.Save(trx)
	if err != nil {
		return Transaction{}, err
	}
	trx.ID = id

	trx.DateRecorded = time.Now()

	return trx, nil
}

func (svc *transactionsService) ViewTransaction(uid string) (Transaction, error) {
	// if _, err := svc.auth.Identity(token); err != nil {
	// 	return Transaction{}, err
	// }
	return svc.store.RetrieveByID(uid)
}

func (svc *transactionsService) ListTransactions(offset, limit uint64) (TransactionPage, error) {
	// if _, err := svc.auth.Identity(token); err != nil {
	// 	return TransactionPage{}, err
	// }
	return svc.store.RetrieveAll(offset, limit)
}

func (svc *transactionsService) ListTransactionsByProperty(prop string, offset, limit uint64) (TransactionPage, error) {
	// if _, err := svc.auth.Identity(token); err != nil {
	// 	return TransactionPage{}, err
	// }
	return svc.store.RetrieveByProperty(prop, offset, limit)
}

func (svc *transactionsService) ListTransactionsByMethod(meth string, offset, limit uint64) (TransactionPage, error) {
	// if _, err := svc.auth.Identity(token); err != nil {
	// 	return TransactionPage{}, err
	// }
	return svc.store.RetrieveByMethod(meth, offset, limit)
}

// func (svc *transactionsService) ListTransactionsByMonth(month string, offset, limit uint64) (TransactionPage, error) {
// 	return TransactionPage{}, nil
// }

// func (svc *transactionsService) ListTransactionsByYear(year string, offset, limit uint64) (TransactionPage, error) {
// 	return TransactionPage{}, nil
// }
