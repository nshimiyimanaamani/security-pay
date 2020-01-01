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
	invoice    uint64
}

// NewRepository ...
func NewRepository(invoice uint64, properties []string) payment.Repository {
	return &repoMock{
		txs:        make(map[string]payment.Transaction),
		properties: properties,
		invoice:    invoice,
	}
}

func (repo *repoMock) Save(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "app/mocks/repoMock.Save"

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

func (repo *repoMock) RetrieveProperty(ctx context.Context, code string) (string, error) {
	const op errors.Op = "app/mocks/repoMock.RetrieveProperty"

	for _, id := range repo.properties {
		if id == code {
			return code, nil
		}
	}
	return "", errors.E(op, "property not found", errors.KindNotFound)
}

func (repo *repoMock) OldestInvoice(ctx context.Context, property string) (uint64, error) {
	const op errors.Op = "app/mocks/repoMock.OldestInvoice"

	for _, val := range repo.properties {
		if val == property {
			return repo.invoice, nil
		}
	}
	return 0, errors.E(op, "no invoice found", errors.KindNotFound)
}
