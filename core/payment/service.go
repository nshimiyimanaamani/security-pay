package payment

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service is the api interface to the payment module
type Service interface {
	// Debit initializes payment from an external account
	Debit(ctx context.Context, tx Transaction) (Response, error)

	// Credit  initializes payment to an external account
	Credit(ctx context.Context, tx Transaction) (Response, error)

	// ProcessDebit processes debit callback
	ProcessDebit(ctx context.Context, res Callback) error

	// ProcessCredit processes credit callback
	ProcessCredit(ctx context.Context, res Callback) error
}

// Options simplifies New func signature
type Options struct {
	Idp     identity.Provider
	Backend Client
	Queue   Queue
	Repo    Repository
}
type service struct {
	backend Client
	idp     identity.Provider
	queue   Queue
	repo    Repository
}

// New initializes the payment service
func New(opts *Options) Service {
	return &service{
		queue:   opts.Queue,
		idp:     opts.Idp,
		backend: opts.Backend,
		repo:    opts.Repo,
	}
}

func (svc service) Debit(ctx context.Context, tx Transaction) (Response, error) {
	const op errors.Op = "core/app/payment/service.Debit"

	failed := Response{TxState: "failed"}
	if err := tx.Validate(); err != nil {
		return failed, errors.E(op, err)
	}

	code, err := svc.repo.RetrieveProperty(ctx, tx.Code)
	if err != nil {
		return failed, errors.E(op, err)
	}

	tx.Code = code

	invoice, err := svc.repo.EarliestInvoice(ctx, tx.Code)
	if err != nil {
		return failed, errors.E(op, err)
	}

	if err := invoice.Satisfy(tx.Amount); err != nil {
		return failed, errors.E(op, err)
	}

	func() {
		tx.Invoice = invoice.ID
		tx.ID = svc.idp.ID()
	}()

	status, err := svc.backend.Pull(ctx, tx)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, tx); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) Credit(ctx context.Context, tx Transaction) (Response, error) {
	const op errors.Op = "core/app/payment/service.Credit"

	failed := Response{TxState: "failed"}

	if err := tx.HackyValidation(); err != nil {
		return failed, errors.E(op, err)
	}
	status, err := svc.backend.Push(ctx, tx)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, tx); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) ProcessDebit(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/app/payment/service.ProcessDebit"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if cb.Data.State != Successful {
		return errors.E(op, "transaction failed unexpectedly", errors.KindUnexpected)
	}

	tx, err := svc.queue.Get(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err)
	}

	if err := svc.repo.Save(ctx, tx); err != nil {
		return errors.E(op, err)
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, tx.ID); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) ProcessCredit(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/app/payment/service.ProcessCredit"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if cb.Data.State != Successful {
		return errors.E(op, "transaction failed unexpectedly", errors.KindUnexpected)
	}
	tx, err := svc.queue.Get(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err)
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, tx.ID); err != nil {
		return errors.E(op, err)
	}
	return nil
}
