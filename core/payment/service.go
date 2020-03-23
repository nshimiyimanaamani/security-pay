package payment

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service is the api interface to the payment module
type Service interface {
	// Initialized by the client app
	Initilize(ctx context.Context, tx Transaction) (Status, error)

	// Validattion is
	Confirm(ctx context.Context, res Callback) error
}

// Options simplifies New func signature
type Options struct {
	Idp     identity.Provider
	Backend Backend
	Queue   Queue
	Repo    Repository
}
type service struct {
	backend Backend
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

func (svc service) Initilize(ctx context.Context, tx Transaction) (Status, error) {
	const op errors.Op = "app.payment.Initialize"

	failed := Status{TxState: "failed"}
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

	//log.Printf("code:%s invoice:%2.f vs entered:%2f", tx.Code, invoice.Amount, tx.Amount)

	if invoice.Amount != tx.Amount {
		return failed, errors.E(op, "amount doesn't match invoice", errors.KindBadRequest)
	}

	identity := func() {
		tx.Invoice = invoice.ID
		tx.ID = svc.idp.ID()
	}
	identity()

	status, err := svc.backend.Pull(ctx, tx)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, tx); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc service) Confirm(ctx context.Context, cb Callback) error {
	const op errors.Op = "app.payment.Confirm"

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
