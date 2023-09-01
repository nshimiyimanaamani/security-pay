package mocks

import (
	"context"

	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/core/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type paymentMock struct{}

// NewPaymentService ...
func NewPaymentService() payment.Service {
	return &paymentMock{}
}

func (svc *paymentMock) Pull(ctx context.Context, tx *payment.TxRequest) (*payment.TxResponse, error) {
	const op errors.Op = "core/ussd/mocks/paymentMock.Initilialize"
	return &payment.TxResponse{
		Status:  "transaction is done",
		TxID:    uuid.New().ID(),
		TxState: "pending",
	}, nil
}

// Validattion is
func (svc *paymentMock) ProcessHook(ctx context.Context, res payment.Callback) error {
	const op errors.Op = "core/ussd/mocks/paymentMock.Confirm"

	return errors.E(op, errors.KindNotImplemented)
}

func (svc *paymentMock) Push(ctx context.Context, tx *payment.TxRequest) (*payment.TxResponse, error) {
	const op errors.Op = "core/ussd/mocks/paymentMock.Initilialize"
	return &payment.TxResponse{
		Status:  "transaction is done",
		TxID:    uuid.New().ID(),
		TxState: "pending",
	}, nil
}

// Validattion is
func (svc *paymentMock) ConfirmPush(ctx context.Context, res payment.Callback) error {
	const op errors.Op = "core/ussd/mocks/paymentMock.Confirm"

	return errors.E(op, errors.KindNotImplemented)
}

func (svc *paymentMock) BulkPull(ctx context.Context, tx *payment.TxRequest, month int) (*payment.TxResponse, error) {
	const op errors.Op = "core/ussd/mocks/paymentMock.BulkPull"
	return nil, errors.E(op, errors.KindNotImplemented)
}

func (svc *paymentMock) CreditPull(ctx context.Context, tx *payment.TxRequest, inv []invoices.Invoice) (*payment.TxResponse, error) {
	const op errors.Op = "core/ussd/mocks/paymentMock.CreditPull"
	return nil, errors.E(op, errors.KindNotImplemented)
}
