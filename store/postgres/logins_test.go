package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginRetrieve(t *testing.T) {
	repo := postgres.NewAuthRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := users.Agent{Telephone: "0780456000", Password: "password", Role: users.Min, Account: account.ID}
	saved, err := saveAgent(t, db, user)

	const op errors.Op = "store/postgres/authRepository.Retrieve"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing user",
			id:   saved.Telephone,
			err:  nil,
		},
		{
			desc: "retrieve non existing user",
			id:   "invalid",
			err:  errors.E(op, "user not found: invalid username or password", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Retrieve(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
