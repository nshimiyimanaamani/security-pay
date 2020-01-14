package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveSectorPayRatio(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
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

	const op errors.Op = "store/postgres/statsRepository.RetrieveSectorPayRatio"

	cases := []struct {
		desc    string
		sector  string
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing sector",
			sector:  property.Address.Sector,
			label:   property.Address.Sector,
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing sector",
			sector:  "invalid",
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "sector not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.RetrieveSectorPayRatio(ctx, tc.sector)
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%s' got '%s'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed'%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected pending '%d' got '%d'\n", tc.desc, tc.pending, pending))
	}
}

func TestRetrieveCellPayRatio(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
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

	const op errors.Op = "store/postgres/statsRepository.RetrieveCellPayRatio"

	cases := []struct {
		desc    string
		cell    string
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing cell",
			cell:    property.Address.Cell,
			label:   property.Address.Cell,
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing cell",
			cell:    "invalid",
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "cell not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.RetrieveCellPayRatio(ctx, tc.cell)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%v' got '%v'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, pending))
	}
}

func TestRetrieveVillagePayRatio(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
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

	const op errors.Op = "store/postgres/statsRepository.RetrieveVillagePayRatio"

	cases := []struct {
		desc    string
		village string
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing village",
			village: property.Address.Village,
			label:   property.Address.Village,
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing village",
			village: "invalid",
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "village not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.RetrieveVillagePayRatio(ctx, tc.village)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%v' got '%v'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, pending))
	}
}
