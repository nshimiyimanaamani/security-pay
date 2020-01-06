package accounts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/accounts/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newService() accounts.Service {
	repo := mocks.NewRepository()
	idp := mocks.NewIdentityProvider()
	opts := &accounts.Options{Repository: repo, IDP: idp}
	return accounts.New(opts)
}

func TestCreate(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/accounts/service.Create"

	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{

		{
			desc:    "create a valid account",
			account: accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     nil,
		},
		{
			desc:    "update account with missing name",
			account: accounts.Account{ID: "paypack.developers", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing name"),
		},
		{
			desc:    "update account with missing type",
			account: accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10},
			err:     errors.E(op, "invalid account: missing type"),
		},
		{
			desc:    "update account with missing account id(sector)",
			account: accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing id"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Create(ctx, tc.account)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/accounts/service.Update"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	ctx := context.Background()

	saved, err := svc.Create(ctx, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{

		{
			desc:    "update account with existing id",
			account: accounts.Account{ID: saved.ID, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     nil,
		},
		{
			desc:    "update account with non existant id",
			account: accounts.Account{ID: "invalid", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "account not found"),
		},
		{
			desc:    "update account with missing name",
			account: accounts.Account{ID: saved.ID, NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing name"),
		},
		{
			desc:    "update account with missing type",
			account: accounts.Account{ID: saved.ID, Name: "remera", NumberOfSeats: 10},
			err:     errors.E(op, "invalid account: missing type"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.Update(ctx, tc.account)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieve(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/accounts/service.Retrieve"

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
	ctx := context.Background()

	saved, err := svc.Create(ctx, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve account with existing id",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "retrieve account with non existant id",
			id:   "invalid",
			err:  errors.E(op, "account not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Retrieve(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestList(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/accounts/service.List"

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.Create(ctx, account)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := []struct {
		desc   string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all accounts",
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the accounts",
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			offset: n + 1,
			limit:  n,
			size:   0,
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
		size := uint64(len(page.Accounts))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}
