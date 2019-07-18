package transactions_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	wrong = "wrong"
)

var transaction = transactions.Transaction{
	ID:       "1000-4433-3343",
	Amount:   "1000.00",
	Method:   "BK",
	Property: "1000-4433-3343",
}

func newService() transactions.Service {
	idp := mocks.NewIdentityProvider()
	store := mocks.NewTransactionStore()
	return transactions.New(idp, store)
}

func TestRecordTransaction(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc        string
		transaction transactions.Transaction
		err         error
	}{
		{
			desc:        "add valid transaction",
			transaction: transaction,
			err:         nil,
		},
		{
			desc:        "add invalid transaction",
			transaction: transactions.Transaction{Amount: "1000.00"},
			err:         transactions.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		_, err := svc.RecordTransaction(tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewTransaction(t *testing.T) {
	svc := newService()
	transaction, _ := svc.RecordTransaction(transaction)

	cases := []struct {
		desc     string
		identity string
		err      error
	}{
		{
			desc:     "view existing transaction",
			identity: transaction.ID,
			err:      nil,
		},
		{
			desc:     "view non-existing transaction",
			identity: wrong,
			err:      transactions.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewTransaction(tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListTransactions(t *testing.T) {
	svc := newService()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.RecordTransaction(transaction)
	}

	cases := []struct {
		desc   string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc: "	list empty set",
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list all transactions",
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half",
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc:   "list last transaction",
			offset: n - 1,
			limit:  n,
			size:   1,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListTransactions(tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByProperty(t *testing.T) {
	svc := newService()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction property for half the transactiona
		if i >= uint64(5) {
			transaction.Property = "1000-4433-0000"
		}
		svc.RecordTransaction(transaction)
	}

	cases := []struct {
		desc     string
		property string
		offset   uint64
		limit    uint64
		size     uint64
		err      error
	}{
		{
			desc:     "list all transactions for an existing property",
			property: transaction.Property,
			offset:   0,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list the last transaction for an existing property",
			property: transaction.Property,
			offset:   n - 1,
			limit:    n,
			size:     1,
			err:      nil,
		},

		{
			desc:     "list half the transaction for an existing property",
			property: transaction.Property,
			offset:   n / 2,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list with zero limit",
			property: transaction.Property,
			offset:   1,
			limit:    0,
			size:     0,
			err:      nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListTransactionsByProperty(tc.property, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByMethod(t *testing.T) {}

func TestListTransactionsByYear(t *testing.T) {}
