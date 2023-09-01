package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveManager(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}
	account = saveAccount(t, db, account)

	const op errors.Op = "store/postgres.userRepository.SaveManager"

	cases := []struct {
		desc string
		user users.Manager
		err  error
	}{
		{
			desc: "save valid manager",
			user: users.Manager{Email: "email@gmail.com", Password: "password", Role: users.Basic, Account: account.ID},
			err:  nil,
		},
		{
			desc: "save duplicate manager",
			user: users.Manager{Email: "email@gmail.com", Password: "password", Role: users.Basic, Account: account.ID},
			err:  errors.E(op, "user already exists", errors.KindAlreadyExists),
		},
		{
			desc: "save manager with invalid data",
			user: users.Manager{Email: "invalid_account@gmail.com", Password: "password", Role: users.Basic, Account: "invalid"},
			err:  errors.E(op, "invalid input data: account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.SaveManager(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieveManager(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}
	account = saveAccount(t, db, account)

	user := users.Manager{Account: account.ID, Email: "email@example.com", Role: users.Basic}
	saved, err := repo.SaveManager(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/userRepository.RetrieveManager"

	cases := []struct {
		desc string
		id   string
		err  error
	}{

		{
			desc: "retrieve existing manager(user)",
			id:   saved.Email,
			err:  nil,
		},
		{
			desc: "retrieve non existing manager(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveManager(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestUpdateManagerCreds(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}
	account = saveAccount(t, db, account)

	user := users.Manager{Account: account.ID, Email: "email@example.com", Role: users.Basic}
	saved, err := repo.SaveManager(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres.userRepository.UpdateManagerCreds"

	cases := []struct {
		desc string
		user users.Manager
		err  error
	}{
		{
			desc: "update existing developer's credentials",
			user: users.Manager{Email: saved.Email, Password: "password"},
			err:  nil,
		},
		{
			desc: "update non existing developer's credentials",
			user: users.Manager{Email: "email2@gmail.com", Password: "password"},
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.UpdateManagerCreds(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListManagers(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}

	account = saveAccount(t, db, account)

	user := users.Manager{Account: account.ID, Role: users.Basic}

	creds := &auth.Credentials{Account: account.ID}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		user.Email = fmt.Sprintf("email%d@gmail.com", i)
		_, err := repo.SaveManager(ctx, user)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "store/postgres.userRepository.ListManagers"

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
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := repo.ListManagers(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Managers))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestDeleteManager(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Bens}
	account = saveAccount(t, db, account)

	user := users.Manager{Account: account.ID, Email: "email@example.com", Role: users.Basic}
	saved, err := repo.SaveManager(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/userRepository.DeleteManager"

	cases := []struct {
		desc string
		id   string
		err  error
	}{

		{
			desc: "retrieve existing manager(user)",
			id:   saved.Email,
			err:  nil,
		},
		{
			desc: "retrieve non existing manager(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.DeleteManager(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}
