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

func TestSaveDeveloper(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	const op errors.Op = "store/postgres.userRepository.SaveDeveloper"

	cases := []struct {
		desc string
		user users.Developer
		err  error
	}{
		{
			desc: "save valid developer",
			user: users.Developer{Email: "email@gmail.com", Password: "password", Role: users.Dev, Account: account.ID},
			err:  nil,
		},
		{
			desc: "save duplicate developer",
			user: users.Developer{Email: "email@gmail.com", Password: "password", Role: users.Dev, Account: account.ID},
			err:  errors.E(op, "user already exists", errors.KindAlreadyExists),
		},
		{
			desc: "save developer with invalid data",
			user: users.Developer{Email: "invalid_account@gmail.com", Password: "password", Role: users.Dev, Account: "invalid"},
			err:  errors.E(op, "invalid input data: account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.SaveDeveloper(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestRetrieveDeveloper(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	user := users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin}
	saved, err := repo.SaveAdmin(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/userRepository.RetrieveDeveloper"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing developer(user)",
			id:   saved.Email,
			err:  nil,
		},
		{
			desc: "retrieve non existing developer(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveDeveloper(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdateDeveloperCreds(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	user := users.Administrator{Account: account.ID, Email: "developer@gmail.com", Role: users.Admin}
	saved, err := repo.SaveAdmin(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres.userRepository.UpdateDeveloperCreds"

	cases := []struct {
		desc string
		user users.Developer
		err  error
	}{
		{
			desc: "update existing developer's credentials",
			user: users.Developer{Email: saved.Email, Password: "password"},
			err:  nil,
		},
		{
			desc: "update non existing developer's credentials",
			user: users.Developer{Email: "email2@gmail.com", Password: "password"},
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.UpdateDeveloperCreds(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListDevelopers(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	user := users.Developer{Account: account.ID, Role: users.Dev}

	creds := &auth.Credentials{Account: account.ID}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		user.Email = fmt.Sprintf("email%d@gmail.com", i)
		_, err := repo.SaveDeveloper(ctx, user)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	const op errors.Op = "store/postgres.userRepository.ListDevelopers"

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{

		"retrieve all developers": {
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of all developers": {
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := repo.ListDevelopers(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Developers))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestDeleteDeveloper(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	user := users.Administrator{Account: account.ID, Email: "email@example.com", Role: users.Admin}
	saved, err := repo.SaveAdmin(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/userRepository.DeleteDeveloper"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing developer(user)",
			id:   saved.Email,
			err:  nil,
		},
		{
			desc: "retrieve non existing developer(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.DeleteDeveloper(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
