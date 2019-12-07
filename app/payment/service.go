package payment

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/identity"
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

	empty := Status{}
	if err := tx.Validate(); err != nil {
		return empty, errors.E(op, err, errors.Kind(err))
	}

	code, err := svc.repo.RetrieveCode(ctx, tx.Code)
	if err != nil {
		return empty, errors.E(op, err, errors.Kind(err))
	}

	tx.ID = svc.idp.ID()
	tx.Code = code

	status, err := svc.backend.Pull(ctx, tx)
	if err != nil {
		return empty, errors.E(op, err, errors.Kind(err))
	}
	if err := svc.queue.Set(ctx, tx); err != nil {
		return empty, errors.E(op, err, errors.Kind(err))
	}
	return status, nil
}

func (svc service) Confirm(ctx context.Context, cb Callback) error {
	const op errors.Op = "app.payment.Confirm"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}

	tx, err := svc.queue.Get(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err, errors.Kind(err))
	}

	if err := svc.repo.Save(ctx, tx); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, tx.ID); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}
	return nil
}
