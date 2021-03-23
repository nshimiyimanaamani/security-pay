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
	payments map[string]payment.Payment
}

// NewPaymentRepository creates an in memory mock of payment.Repository
func NewPaymentRepository() payment.Repository {
	return &repositoryMock{
		payments: make(map[string]payment.Payment),
	}
}

func (repo *repositoryMock) Save(ctx context.Context, payment payment.Payment) error {
	const op errors.Op = "core/payment/mocks/repositoryMock.Save"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, py := range repo.payments {
		if py.ID == payment.ID {
			return errors.E(op, "payment not found", errors.KindNotFound)
		}
	}

	repo.counter++
	repo.payments[payment.ID] = payment

	return nil
}

func (repo *repositoryMock) Find(ctx context.Context, id string) (payment.Payment, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.Find"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	val, ok := repo.payments[id]
	if !ok {
		return payment.Payment{}, errors.E(op, "payment not found", errors.KindNotFound)
	}
	return val, nil
}

func (repo *repositoryMock) Update(ctx context.Context, payment payment.Payment) error {
	const op errors.Op = "core/payment/mocks/repositoryMock.Update"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.payments[payment.ID]; !ok {
		return errors.E(op, "payment not found", errors.KindNotFound)
	}

	repo.payments[payment.ID] = payment
	return nil
}
