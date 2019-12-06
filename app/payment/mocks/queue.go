package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Queue) = (*mockQueue)(nil)

// mock is a payment.Queue implementation that doesn't support concurency
type mockQueue struct {
	transactions map[string]interface{}
}

// NewQueue initialises the mockQueue
func NewQueue() payment.Queue {
	return &mockQueue{
		transactions: make(map[string]interface{}),
	}
}

func (q *mockQueue) Set(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "mocksQueue.Set"

	return errors.E(op, errors.KindNotFound)
}

func (q *mockQueue) Unset(ctx context.Context, uid string) (payment.Transaction, error) {
	const op errors.Op = "mocksQueue.Unset"

	return payment.Transaction{}, errors.E(op, errors.KindNotFound)
}
