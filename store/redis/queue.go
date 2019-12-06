package bigcache

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Queue) = (*queue)(nil)

type queue struct{}

// NewQueue initialises the queue
func NewQueue(url string) payment.Queue {
	return &queue{}
}

func (q *queue) Set(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "queue.Set"

	return errors.E(op, errors.KindNotFound)
}

func (q *queue) Unset(ctx context.Context, uid string) (payment.Transaction, error) {
	const op errors.Op = "queue.Unset"

	return payment.Transaction{}, errors.E(op, errors.KindNotFound)
}
