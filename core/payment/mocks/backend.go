package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Client) = (*backendMock)(nil)

type backendMock struct{}

// NewBackend instantiates a payment backend mock
func NewBackend() payment.Client {
	return &backendMock{}
}

func (bc *backendMock) Pull(ctx context.Context, tx *payment.TxRequest) (*payment.TxResponse, error) {
	const op errors.Op = "core/payment/mocks/backendMock.Pull"

	return &payment.TxResponse{
		TxID:    tx.ID,
		Status:  "success",
		TxState: "processing",
	}, nil
}

func (bc *backendMock) Push(ctx context.Context, tx *payment.TxRequest) (*payment.TxResponse, error) {
	const op errors.Op = "core/payment/mocks/backendMock.Push"

	return &payment.TxResponse{
		TxID:    tx.ID,
		Status:  "success",
		TxState: "processing",
	}, nil
}
