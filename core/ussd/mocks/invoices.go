package mocks

// import (
// 	"context"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/rugwirobaker/paypack-backend/core/payment"
// 	"github.com/rugwirobaker/paypack-backend/pkg/errors"
// )

// type invoicesMock struct {
// 	counter    uint64
// 	mu         sync.Mutex
// 	txs        map[string]payment.Transaction
// 	properties []string
// 	invoice    payment.Invoice
// }

// // NewInvoicesRepository ...
// func NewInvoicesRepository(inv payment.Invoice, properties []string) payment.Repository {
// 	return &invoicesMock{
// 		txs:        make(map[string]payment.Transaction),
// 		properties: properties,
// 		invoice:    inv,
// 	}
// }

// func (repo *invoicesMock) Save(ctx context.Context, tx payment.Transaction) error {
// 	const op errors.Op = "app/mocks/repoMock.Save"

// 	repo.mu.Lock()
// 	defer repo.mu.Unlock()

// 	for _, saved := range repo.txs {
// 		if saved.ID == tx.ID {
// 			return errors.E(op, "transaction already exist", errors.KindAlreadyExists)
// 		}
// 	}

// 	tx.RecordedAt = time.Now()

// 	repo.counter++
// 	tx.ID = strconv.FormatUint(repo.counter, 10)
// 	repo.txs[tx.ID] = tx
// 	return nil
// }

// func (repo *invoicesMock) RetrieveProperty(ctx context.Context, code string) (string, error) {
// 	const op errors.Op = "app/mocks/repoMock.RetrieveProperty"

// 	for _, id := range repo.properties {
// 		if id == code {
// 			return code, nil
// 		}
// 	}
// 	return "", errors.E(op, "property not found", errors.KindNotFound)
// }

// func (repo *invoicesMock) EarliestInvoice(ctx context.Context, property string) (payment.Invoice, error) {
// 	const op errors.Op = "app/mocks/repoMock.OldestInvoice"

// 	for _, val := range repo.properties {
// 		if val == property {
// 			return repo.invoice, nil
// 		}
// 	}
// 	return payment.Invoice{}, errors.E(op, "no invoice found", errors.KindNotFound)
// }
