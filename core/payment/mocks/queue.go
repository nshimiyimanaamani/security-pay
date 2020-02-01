package mocks

import (
	"context"
	"sync"

	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Queue) = (*mockQueue)(nil)

// mock is a payment.Queue implementation that doesn't support concurency
type mockQueue struct {
	mu  sync.Mutex
	txs map[string]payment.Transaction
}

// NewQueue initialises the mockQueue
func NewQueue() payment.Queue {
	return &mockQueue{
		txs: make(map[string]payment.Transaction),
	}
}

func (q *mockQueue) Set(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "mocksQueue.Set"

	q.mu.Lock()
	defer q.mu.Unlock()

	for _, saved := range q.txs {
		if saved.ID == tx.ID {
			return errors.E(op, "transaction already exist", errors.KindAlreadyExists)
		}
	}
	q.txs[tx.ID] = tx
	return nil
}

func (q *mockQueue) Get(ctx context.Context, uid string) (payment.Transaction, error) {
	const op errors.Op = "mocksQueue.Get"

	q.mu.Lock()
	defer q.mu.Unlock()

	val, ok := q.txs[uid]
	if !ok {
		return payment.Transaction{}, errors.E(op, "transaction not found", errors.KindNotFound)
	}

	return val, nil
}

func (q *mockQueue) Remove(ctx context.Context, uid string) error {
	const op errors.Op = "mocksQueue.Remove"

	q.mu.Lock()
	defer q.mu.Unlock()

	_, ok := q.txs[uid]
	if !ok {
		return errors.E(op, "transaction not found", errors.KindNotFound)
	}

	delete(q.txs, uid)
	return nil
}
