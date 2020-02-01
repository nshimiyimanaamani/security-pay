package transactions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/core/transactions/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var transaction = transactions.Transaction{
	ID:      "1000-4433-3343",
	Amount:  1000.00,
	Method:  "bk",
	MadeFor: "1000-4433-3343",
	OwnerID: "1000-4433-3343",
}

func newService() transactions.Service {
	repo := mocks.NewRepository()
	idp := mocks.NewIdentityProvider()
	opts := &transactions.Options{
		Repo: repo,
		Idp:  idp,
	}
	return transactions.New(opts)
}

func TestRecordTransaction(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/transactions/service.Record"

	cases := []struct {
		desc        string
		token       string
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
			transaction: transactions.Transaction{Amount: 1000.00},
			err:         errors.E(op, "invalid transaction"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Record(ctx, tc.transaction)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieveTransaction(t *testing.T) {
	svc := newService()

	ctx := context.Background()
	transaction, err := svc.Record(ctx, transaction)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "app/transactions/service.Retrieve"

	cases := []struct {
		desc     string
		token    string
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
			err:      errors.E(op, "transaction not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Retrieve(ctx, tc.identity)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListTransactions(t *testing.T) {
	svc := newService()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	}

	const op errors.Op = "app/transactions/service.List"

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
		ctx := context.Background()
		page, err := svc.List(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByProperty(t *testing.T) {
	svc := newService()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction property for half the transactiona
		if i >= uint64(5) {
			transaction.MadeFor = "1000-4433-0000"
		}
		ctx := context.Background()
		_, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "app/transactions/service.ListByProperty"

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
			property: transaction.MadeFor,
			offset:   0,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list the last transaction for an existing property",
			property: transaction.MadeFor,
			offset:   n - 1,
			limit:    n,
			size:     1,
			err:      nil,
		},

		{
			desc:     "list half the transaction for an existing property",
			property: transaction.MadeFor,
			offset:   n / 2,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
		{
			desc:     "list with zero limit",
			property: transaction.MadeFor,
			offset:   1,
			limit:    0,
			size:     0,
			err:      nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListByProperty(ctx, tc.property, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByPropertyR(t *testing.T) {
	svc := newService()

	t.Skip()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction property for half the transactiona
		if i >= uint64(5) {
			transaction.MadeFor = "1000-4433-0000"
		}
		ctx := context.Background()
		_, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "app/transactions/service.ListByProperty"

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
			property: transaction.MadeFor,
			offset:   0,
			limit:    n,
			size:     n / 2,
			err:      nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListByPropertyR(ctx, tc.property)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestListTransactionsByMethod(t *testing.T) {
	svc := newService()

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		//change transaction method for half the transactions
		if i >= uint64(5) {
			transaction.Method = "mtn"
		}
		ctx := context.Background()
		_, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "app/transactions/service.ListByMethod"

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
			desc:   "list all transactions with existing transaction method",
			method: transaction.Method,
			offset: 0,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc:   "list the last transaction with existant transaction method",
			method: transaction.Method,
			offset: n - 1,
			limit:  n,
			size:   1,
			err:    nil,
		},

		{
			desc:   "list half the transaction with existant transaction method",
			method: transaction.Method,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			method: transaction.Method,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListByMethod(ctx, tc.method, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
