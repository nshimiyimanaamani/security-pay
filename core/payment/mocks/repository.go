package mocks

import (
	"context"
	"sync"

	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type repositoryMock struct {
	mu       sync.Mutex
	counter  uint64
	payments map[string]payment.TxRequest
}

// NewPaymentRepository creates an in memory mock of payment.Repository
func NewPaymentRepository() payment.Repository {
	return &repositoryMock{
		payments: make(map[string]payment.TxRequest),
	}
}

func (repo *repositoryMock) Save(ctx context.Context, payment *payment.TxRequest) error {
	const op errors.Op = "core/payment/mocks/repositoryMock.Save"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, py := range repo.payments {
		if py.ID == payment.ID {
			return errors.E(op, "payment not found", errors.KindNotFound)
		}
	}

	repo.counter++
	repo.payments[payment.ID] = *payment

	return nil
}

func (repo *repositoryMock) Find(ctx context.Context, id string) ([]*payment.TxRequest, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.Find"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return nil, errors.E(op, "not implemented", errors.KindUnexpected)
}

func (repo *repositoryMock) Update(ctx context.Context, status string, payment []*payment.TxRequest) error {
	const op errors.Op = "core/payment/mocks/repositoryMock.Update"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, "not implemented", errors.KindUnexpected)
}

func (repo *repositoryMock) BulkSave(ctx context.Context, payment []*payment.TxRequest) error {
	const op errors.Op = "core/payment/mocks/repositoryMock.Update"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, "not implemented", errors.KindUnexpected)
}
