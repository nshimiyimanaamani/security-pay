package mocks

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type repoMock struct {
	counter    uint64
	mu         sync.Mutex
	txs        map[string]payment.Transaction
	properties []string
}

// NewRepository ...
func NewRepository(properties []string) payment.Repository {
	return &repoMock{
		txs:        make(map[string]payment.Transaction),
		properties: properties,
	}
}

func (repo *repoMock) Save(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "mocks.repo.Save"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, saved := range repo.txs {
		if saved.ID == tx.ID {
			return errors.E(op, "transaction already exist", errors.KindAlreadyExists)
		}
	}

	tx.RecordedAt = time.Now()

	repo.counter++
	tx.ID = strconv.FormatUint(repo.counter, 10)
	repo.txs[tx.ID] = tx
	return nil
}

func (repo *repoMock) RetrieveCode(ctx context.Context, code string) (string, error) {
	const op errors.Op = "mocks.repo.RetrieveCode"

	for _, id := range repo.properties {
		if id == code {
			return code, nil
		}
	}
	return "", errors.E(op, "property not found", errors.KindNotFound)
}
