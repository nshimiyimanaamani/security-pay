package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Backend) = (*backendMock)(nil)

type backendMock struct{}

// NewBackend instantiates a payment backend mock
func NewBackend() payment.Backend {
	return &backendMock{}
}

func (bc *backendMock) Pull(ctx context.Context, tx payment.Transaction) (payment.Status, error) {
	const op errors.Op = "backendMock.Pull"

	return payment.Status{
		TxID:    tx.ID,
		Status:  "success",
		TxState: "processing",
	}, nil
}

func (bc *backendMock) Status(context.Context) (int, error) {
	const op errors.Op = "backendMock.Status"

	return 0, errors.E(op, errors.KindNotImplemented)
}

func (bc *backendMock) Auth(appID, appSecret string) (string, error) {
	const op errors.Op = "backendMock.Auth"

	return "", errors.E(op, errors.KindNotImplemented)
}
