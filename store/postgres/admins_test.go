package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAdmin(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	const op errors.Op = "store/postgres/userRepository.SaveAdmin"

	id := "gasabo.remera"

	account := accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		user users.Administrator
		err  error
	}{
		{
			desc: "save valid administrator",
			user: users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin},
			err:  nil,
		},
		{
			desc: "save duplicate administrator",
			user: users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin},
			err:  errors.E(op, "user already exists", errors.KindAlreadyExists),
		},
		{
			desc: "save administrator with invalid account",
			user: users.Administrator{Account: "invalid", Email: "email2@example.com", Role: users.Admin},
			err:  errors.E(op, "invalid input data: account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.SaveAdmin(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieveAdmin(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	id := "gasabo.remera"

	account := accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	user := users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin}
	saved, err := repo.SaveAdmin(context.Background(), user)

	const op errors.Op = "store/postgres/userRepository.RetrieveAdmin"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing admin(user)",
			id:   saved.Email,
			err:  nil,
		},
		{
			desc: "retrieve non existing admin(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveAdmin(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdateAdminCreds(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	id := "gasabo.remera"
	account := accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	user := users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin}
	saved, err := repo.SaveAdmin(context.Background(), user)

	const op errors.Op = "store/postgres/userRepository.UpdateAdminCreds"

	cases := []struct {
		desc string
		user users.Administrator
		err  error
	}{
		{
			desc: "update existing admin's credentials",
			user: users.Administrator{Email: saved.Email, Password: "password"},
			err:  nil,
		},
		{
			desc: "update non existing admin's credentials",
			user: users.Administrator{Email: "email2@gmail.com", Password: "password"},
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.UpdateAdminCreds(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListAdmins(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	user := users.Administrator{Account: account.ID, Role: users.Admin}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		user.Email = fmt.Sprintf("email%d@gmail.com", i)
		_, err := repo.SaveAdmin(ctx, user)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "store/postgres.userRepository.ListAdmins"

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
		page, err := repo.ListAdmins(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Administrators))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}
