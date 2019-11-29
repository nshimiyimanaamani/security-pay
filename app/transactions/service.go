package transactions

import (
	"context"
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
	// Record adds a new transaction record
	Record(ctx context.Context, tx Transaction) (Transaction, error)

	// Retrieve retrieves data about the transaction identified with the provided ID,
	Retrieve(ctx context.Context, uid string) (Transaction, error)

	// List retrieves data about subset of transactions that belongs to the
	// user identified by the provided key.
	List(ctx context.Context, offset, limit uint64) (TransactionPage, error)

	// ListByProperty retrieves data about a subset of transactions that belong to
	// a property identified by the given id. i.e filer by property
	ListByProperty(ctx context.Context, p string, offset, limit uint64) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that were made using
	// a given method.
	ListByPeriod(ctx context.Context, offset, limit uint64) (TransactionPage, error)
}

var _ Service = (*transactionsService)(nil)

type transactionsService struct {
	idp  identity.Provider
	repo Repository
}

// Options ...
type Options struct {
	Idp  identity.Provider
	Repo Repository
}

//New instantiates a new transaxtions service
func New(opts *Options) Service {
	return &transactionsService{
		idp:  opts.Idp,
		repo: opts.Repo,
	}
}

func (svc *transactionsService) Record(ctx context.Context, tx Transaction) (Transaction, error) {
	if err := tx.Validate(); err != nil {
		return Transaction{}, err
	}

	tx.ID = svc.idp.ID()

	id, err := svc.repo.Save(ctx, tx)
	if err != nil {
		return Transaction{}, err
	}
	tx.ID = id

	tx.DateRecorded = time.Now()

	return tx, nil
}

func (svc *transactionsService) Retrieve(ctx context.Context, uid string) (Transaction, error) {
	return svc.repo.RetrieveByID(ctx, uid)
}

func (svc *transactionsService) List(ctx context.Context, offset, limit uint64) (TransactionPage, error) {
	return svc.repo.RetrieveAll(ctx, offset, limit)
}

func (svc *transactionsService) ListByProperty(ctx context.Context, p string, offset, limit uint64) (TransactionPage, error) {
	return svc.repo.RetrieveByProperty(ctx, p, offset, limit)
}

func (svc *transactionsService) ListByPeriod(ctx context.Context, offset, limit uint64) (TransactionPage, error) {
	return svc.repo.RetrieveByPeriod(ctx, offset, limit)
}
