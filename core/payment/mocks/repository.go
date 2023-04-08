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

func (repo *repositoryMock) List(ctx context.Context, flts *payment.Filters) (payment.PaymentResponse, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.List"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return payment.PaymentResponse{}, errors.E(op, "not implemented", errors.KindUnexpected)
}

func (repo *repositoryMock) ListDailyTransactions(ctx context.Context, flts *payment.MetricFilters) (payment.Transactions, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.ListDailyTransactions"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return payment.Transactions{}, errors.E(op, "not implemented", errors.KindUnexpected)
}

func (repo *repositoryMock) TodayTransaction(ctx context.Context, flts *payment.MetricFilters) (payment.Transaction, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.TodayTransaction"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return payment.Transaction{}, errors.E(op, "not implemented", errors.KindUnexpected)
}

func (repo *repositoryMock) TodaySummary(ctx context.Context, flts *payment.MetricFilters) (payment.Summaries, error) {
	const op errors.Op = "core/payment/mocks/repositoryMock.TodaySummary"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return payment.Summaries{}, errors.E(op, "not implemented", errors.KindUnexpected)
}
