package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAgent(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	phone := "0780456000"

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	const op errors.Op = "store/postgres/userRepository.SaveAgent"

	cases := []struct {
		desc string
		user users.Agent
		err  error
	}{
		{
			desc: "save valid agent",
			user: users.Agent{Telephone: phone, Password: "password", Role: users.Min, Account: account.ID},
			err:  nil,
		},
		{
			desc: "save duplicate agent",
			user: users.Agent{Telephone: phone, Password: "password", Role: users.Min, Account: account.ID},
			err:  errors.E(op, "user already exists", errors.KindAlreadyExists),
		},
		{
			desc: "save agent with invalid data",
			user: users.Agent{Telephone: "0780456450", Password: "password", Role: users.Min, Account: "invalid"},
			err:  errors.E(op, "invalid input data: account not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.SaveAgent(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestRetrieveAgent(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := users.Agent{Telephone: "0780456000", Password: "password", Role: users.Min, Account: account.ID}
	saved, err := repo.SaveAgent(context.Background(), user)

	const op errors.Op = "store/postgres/userRepository.RetrieveAgent"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing agent(user)",
			id:   saved.Telephone,
			err:  nil,
		},
		{
			desc: "retrieve non existing agent(user)",
			id:   "invalid",
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveAgent(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdateAgentDetails(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := users.Agent{Telephone: "0780456000", Password: "password", Role: users.Min, Account: account.ID}
	saved, err := repo.SaveAgent(context.Background(), user)

	const op errors.Op = "store/postgres/userRepository.UpdateAgentDetails"

	cases := []struct {
		desc string
		user users.Agent
		err  error
	}{
		{
			desc: "update existing agent's credentials",
			user: users.Agent{Telephone: saved.Telephone, FirstName: "fname", LastName: "lname", UpdatedAt: time.Now()},
			err:  nil,
		},
		{
			desc: "update non existing agent's credentials",
			user: users.Agent{Telephone: "0781406751", UpdatedAt: time.Now()},
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.UpdateAgentDetails(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdateAgentCreds(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := users.Agent{Telephone: "0780456000", Password: "password", Role: users.Min, Account: account.ID}
	saved, err := repo.SaveAgent(context.Background(), user)

	const op errors.Op = "store/postgres/userRepository.UpdateAgentCreds"

	cases := []struct {
		desc string
		user users.Agent
		err  error
	}{
		{
			desc: "update existing agent's credentials",
			user: users.Agent{Telephone: saved.Telephone, Password: "password", UpdatedAt: time.Now()},
			err:  nil,
		},
		{
			desc: "update non existing agent's credentials",
			user: users.Agent{Telephone: "0781406751", Password: "password", UpdatedAt: time.Now()},
			err:  errors.E(op, "user not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.UpdateAgentCreds(ctx, tc.user)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListAgents(t *testing.T) {
	repo := postgres.NewUserRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := users.Agent{Account: account.ID, Role: users.Min}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		user.Telephone = fmt.Sprintf("email%d@gmail.com", i)
		_, err := repo.SaveAgent(ctx, user)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	const op errors.Op = "store/postgres/userRepository.ListAgents"

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{

		"retrieve all agents": {
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of all agents": {
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.ListAgents(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Agents))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}
