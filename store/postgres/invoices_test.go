package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAll(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	// save agent
	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent, err = saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save owner
	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc     string
		property string
		months   uint
		size     uint
		total    uint
		err      error
	}{
		{
			desc:     "retrieve invoices for existant property",
			property: property.ID,
			months:   1,
			size:     1,
			total:    1,
			err:      nil,
		},
		{
			desc:     "retrieve invoices for non-existant property",
			property: "invalid",
			months:   1,
			size:     0,
			total:    0,
			err:      nil,
		},
	}

	for _, tc := range cases {
		page, err := repo.ListAll(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}

func TestListPending(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	// save agent
	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent, err = saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save owner
	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc     string
		property string
		months   uint
		size     uint
		total    uint
		err      error
	}{
		// {
		// 	desc:     "retrieve invoices for existant property",
		// 	property: property.ID,
		// 	months:   1,
		// 	size:     1,
		// 	total:    1,
		// 	err:      nil,
		// },
		// {
		// 	desc:     "retrieve invoices for non-existant property",
		// 	property: "invalid",
		// 	months:   1,
		// 	size:     0,
		// 	total:    0,
		// 	err:      nil,
		// },
	}

	for _, tc := range cases {
		page, err := repo.ListPending(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}

func TestListPayed(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	// save agent
	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent, err = saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save owner
	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc     string
		property string
		months   uint
		size     uint
		total    uint
		err      error
	}{
		// {
		// 	desc:     "retrieve invoices for existant property",
		// 	property: property.ID,
		// 	months:   1,
		// 	size:     1,
		// 	total:    1,
		// 	err:      nil,
		// },
		// {
		// 	desc:     "retrieve invoices for non-existant property",
		// 	property: "invalid",
		// 	months:   1,
		// 	size:     0,
		// 	total:    0,
		// 	err:      nil,
		// },
	}

	for _, tc := range cases {
		page, err := repo.ListPayed(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}
