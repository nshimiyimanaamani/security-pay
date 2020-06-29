package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type paymentMock struct{}

// NewPaymentService ...
func NewPaymentService() payment.Service {
	return &paymentMock{}
}

func (svc *paymentMock) Initilize(ctx context.Context, tx payment.Transaction) (payment.Status, error) {
	const op errors.Op = "core/ussd/mocks/paymentMock.Initilialize"
	return payment.Status{
		Status: "transaction is done",
	}, nil
}

// Validattion is
func (svc *paymentMock) Confirm(ctx context.Context, res payment.Callback) error {
	const op errors.Op = "core/ussd/mocks/paymentMock.Confirm"

	return errors.E(op, errors.KindNotImplemented)
}
