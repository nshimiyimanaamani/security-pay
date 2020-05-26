package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAccount(t *testing.T) {
	repo := postgres.NewAccountRepository(db)
	defer CleanDB(t, db)

	const op errors.Op = "store/postgres.accountRepository.Save"

	id := "gasabo.gisozi"

	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{
		{
			desc:    "save valid account",
			account: accounts.Account{ID: id, Name: "gisozi", NumberOfSeats: 10, Type: accounts.Bens},
			err:     nil,
		},
		{
			desc:    "save already existing account",
			account: accounts.Account{ID: id, Name: "gisozi", NumberOfSeats: 10, Type: accounts.Bens},
			err:     errors.E(op, "account already exists", errors.KindAlreadyExists),
		},
		// {
		// 	desc:    "save account with invalid account type",
		// 	account: accounts.Account{ID: "gasabo.muhazi", Name: "remera", NumberOfSeats: 10},
		// 	err:     errors.E(op, "invalid account data ", errors.KindBadRequest),
		// },
		// {
		// 	desc:    "save account with invalid id",
		// 	account: accounts.Account{ID: "invalid", Name: "gisozi", NumberOfSeats: 10, Type: accounts.Bens},
		// 	err:     errors.E(op, "invalid input data: sector not found", errors.KindNotFound),
		// },
		// {
		// 	desc:    "save account with empty name",
		// 	account: accounts.Account{ID: "invalid", NumberOfSeats: 10, Type: accounts.Devs},
		// 	err:     errors.E(op, "invalid account data ", errors.KindBadRequest),
		// },

	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Save(ctx, tc.account)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdateAccount(t *testing.T) {
	repo := postgres.NewAccountRepository(db)
	defer CleanDB(t, db)

	const op errors.Op = "store/postgres.accountRepository.Update"

	id := "paypack.developers"

	account := accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	saved, err := repo.Save(context.Background(), account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{
		{
			desc:    "update existant account",
			account: saved,
			err:     nil,
		},
		{
			desc:    "update non-existant account",
			account: accounts.Account{ID: uuid.New().ID(), Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Update(ctx, tc.account)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieveAccount(t *testing.T) {
	repo := postgres.NewAccountRepository(db)
	defer CleanDB(t, db)

	const op errors.Op = "store/postgres.accountRepository.Retrieve"

	id := "paypack.developers"

	account := accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	saved, err := repo.Save(context.Background(), account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing account",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing account",
			id:   uuid.New().ID(),
			err:  errors.E(op, "account not found", errors.KindNotFound),
		},
		{
			desc: "retrieve owner with malformed id",
			id:   "invalid",
			err:  errors.E(op, "account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Retrieve(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListAccounts(t *testing.T) {
	repo := postgres.NewAccountRepository(db)

	defer CleanDB(t, db)

	data := []struct {
		id          string
		name        string
		accountType accounts.AccountType
	}{
		{id: "paypack.test", name: "test", accountType: accounts.Devs},
		{id: "paypack.developers", name: "developers", accountType: accounts.Devs},
		{id: "gasabo.remera", name: "remera", accountType: accounts.Bens},
		{id: "gasabo.kimironko", name: "kimironko", accountType: accounts.Bens},
	}

	n := uint64(4)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()

		account := accounts.Account{ID: data[i].id, Name: data[i].name, NumberOfSeats: 10, Type: data[i].accountType}

		_, err := repo.Save(ctx, account)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all owners": {
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of all owners": {
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.List(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Accounts))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}
