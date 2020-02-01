package transactions

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// var (
// 	// ErrConflict attempt to create an entity with an alreasdy existing id
// 	ErrConflict = errors.New("entity already exists")

// 	//ErrInvalidEntity indicates malformed entity specification (e.g.
// 	//invalid username,  password, account).
// 	ErrInvalidEntity = errors.New("invalid entity format")

// 	// ErrNotFound indicates a non-existent entity request.
// 	ErrNotFound = errors.New("non-existent entity")
// )

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

	// ListByProperty retrieves data about a subset of transactions that belong to
	// a property identified by the given id. i.e filer by property(Forgive for this if u are maintaing)
	//this my standards where higher but I had to lower them to fit in
	ListByPropertyR(ctx context.Context, p string) (TransactionPage, error)

	// ListTransactionByDate retrieves data about a subset of transactions that were made using
	// a given method.
	ListByMethod(ctx context.Context, m string, offset, limit uint64) (TransactionPage, error)
}

var _ Service = (*service)(nil)

type service struct {
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
	return &service{
		idp:  opts.Idp,
		repo: opts.Repo,
	}
}

func (svc *service) Record(ctx context.Context, tx Transaction) (Transaction, error) {
	const op errors.Op = "app/transactions/service.Record"

	if err := tx.Validate(); err != nil {
		return Transaction{}, errors.E(op, err)
	}

	tx.ID = svc.idp.ID()

	id, err := svc.repo.Save(ctx, tx)
	if err != nil {
		return Transaction{}, errors.E(op, err)
	}
	tx.ID = id

	return tx, nil
}

func (svc *service) Retrieve(ctx context.Context, uid string) (Transaction, error) {
	const op errors.Op = "app/transactions/service.Retrieve"

	tx, err := svc.repo.RetrieveByID(ctx, uid)
	if err != nil {
		return Transaction{}, errors.E(op, err)
	}
	return tx, nil
}

func (svc *service) List(ctx context.Context, offset, limit uint64) (TransactionPage, error) {
	const op errors.Op = "app/transactions/service.List"

	page, err := svc.repo.RetrieveAll(ctx, offset, limit)
	if err != nil {
		return TransactionPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) ListByProperty(ctx context.Context, p string, offset, limit uint64) (TransactionPage, error) {
	const op errors.Op = "app/transactions/service.ListByProperty"

	page, err := svc.repo.RetrieveByProperty(ctx, p, offset, limit)
	if err != nil {
		return TransactionPage{}, errors.E(op, err)
	}

	return page, nil
}

func (svc *service) ListByPropertyR(ctx context.Context, p string) (TransactionPage, error) {
	const op errors.Op = "app/transactions/service.ListByProperty"

	page, err := svc.repo.RetrieveByPropertyR(ctx, p)
	if err != nil {
		return TransactionPage{}, errors.E(op, err)
	}

	return page, nil
}

func (svc *service) ListByMethod(ctx context.Context, m string, offset, limit uint64) (TransactionPage, error) {
	const op errors.Op = "app/transactions/service.ListByPeriod"

	page, err := svc.repo.RetrieveByMethod(ctx, m, offset, limit)
	if err != nil {
		return TransactionPage{}, errors.E(op, err)
	}
	return page, nil
}
