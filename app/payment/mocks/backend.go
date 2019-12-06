package mocks

import "github.com/rugwirobaker/paypack-backend/app/payment"

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

var _ (payment.Backend) = (*backendMock)(nil)

type backendMock struct{}

// NewBackend instantiates a payment backend mock
func NewBackend() payment.Backend {
	return &backendMock{}
}

func (bc *backendMock) Pull(ctx context.Context, tx payment.Transaction) (string, error) {
	const op errors.Op = "backendMock.Pull"

	return "", errors.E(op, errors.KindNotImplemented)
}

func (bc *backendMock) Status(context.Context) (int, error) {
	const op errors.Op = "backendMock.Status"

	return 0, errors.E(op, errors.KindNotImplemented)
}

func (bc *backendMock) Auth(ctx context.Context) error {
	const op errors.Op = "backendMock.Auth"

	return errors.E(op, errors.KindNotImplemented)
}
