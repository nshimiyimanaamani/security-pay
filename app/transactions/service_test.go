package transactions_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/users"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions/mocks"
	"github.com/stretchr/testify/assert"
)

var transaction = transactions.Transaction{
	ID:      "1000-4433-3343",
	Amount:  1000.00,
	Method:  "BK",
	MadeFor: "1000-4433-3343",
	MadeBy:  "1000-4433-3343",
}

func newService(tokens map[string]string) transactions.Service {
	auth := mocks.NewAuthBackend(tokens)
	idp := mocks.NewIdentityProvider()
	store := mocks.NewTransactionStore()
	return transactions.New(idp, store, auth)
}

func TestRecordTransaction(t *testing.T) {
	svc := newService(map[string]string{token: email})

	cases := []struct {
		desc        string
		token       string
		transaction transactions.Transaction
		err         error
	}{
		{
			desc:        "add valid transaction",
			token:       token,
			transaction: transaction,
			err:         nil,
		},
		{
			desc:        "add invalid transaction",
			token:       token,
			transaction: transactions.Transaction{Amount: 1000.00},
			err:         transactions.ErrInvalidEntity,
		},
		{
			desc:        "add transaction with invalid token",
			transaction: transaction,
			err:         users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		_, err := svc.RecordTransaction(tc.token, tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewTransaction(t *testing.T) {
	svc := newService(map[string]string{token: email})
	transaction, _ := svc.RecordTransaction(token, transaction)

	cases := []struct {
		desc     string
		token    string
		identity string
		err      error
	}{
		{
			desc:     "view existing transaction",
			token:    token,
			identity: transaction.ID,
			err:      nil,
		},
		{
			desc:     "view non-existing transaction",
			token:    token,
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
	svc := newService(map[string]string{token: email})

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.RecordTransaction(token, transaction)
	}

	cases := []struct {
		desc   string
		token  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc: "	list empty set",
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list all transactions",
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half",
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc:   "list last transaction",
			token:  token,
			offset: n - 1,
			limit:  n,
			size:   1,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			token:  token,
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
	svc := newService(map[string]string{token: email})

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction property for half the transactiona
		if i >= uint64(5) {
			transaction.MadeFor = "1000-4433-0000"
		}
		svc.RecordTransaction(token, transaction)
	}

	cases := []struct {
		desc     string
		token    string
		property string
		offset   uint64
		limit    uint64
		size     uint64
		err      error
	}{
		{
			desc:     "list all transactions for an existing property",
			token:    token,
			property: transaction.MadeFor,
			offset:   0,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list the last transaction for an existing property",
			token:    token,
			property: transaction.MadeFor,
			offset:   n - 1,
			limit:    n,
			size:     1,
			err:      nil,
		},

		{
			desc:     "list half the transaction for an existing property",
			token:    token,
			property: transaction.MadeFor,
			offset:   n / 2,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list with zero limit",
			token:    token,
			property: transaction.MadeFor,
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

func TestListTransactionsByMethod(t *testing.T) {
	svc := newService(map[string]string{token: email})

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction property for half the transactiona
		if i >= uint64(5) {
			transaction.MadeFor = "1000-4433-0000"
		}
		svc.RecordTransaction(token, transaction)
	}

	cases := []struct {
		desc   string
		token  string
		method string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all transactions made with given method",
			token:  token,
			method: transaction.Method,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half the transaction with given method",
			token:  token,
			method: transaction.Method,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			token:  token,
			method: transaction.Method,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListTransactionsByMethod(tc.method, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByMonth(t *testing.T) {}

func TestListTransactionsByYear(t *testing.T) {}
