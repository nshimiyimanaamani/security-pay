package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveSectorPayRatio(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindSectorRatio"

	cases := []struct {
		desc    string
		sector  string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing sector",
			sector:  property.Address.Sector,
			label:   property.Address.Sector,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing sector",
			sector:  "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "no data found for this sector", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindSectorRatio(ctx, tc.sector, tc.year, tc.month)
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

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindCellRatio"

	cases := []struct {
		desc    string
		cell    string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing cell",
			cell:    property.Address.Cell,
			label:   property.Address.Cell,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing cell",
			cell:    "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "no data found for this cell", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindCellRatio(ctx, tc.cell, tc.year, tc.month)
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

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindVillageRatio"

	cases := []struct {
		desc    string
		village string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{
		{
			desc:    "retrieve payment ratio for existing village",
			village: property.Address.Village,
			label:   property.Address.Village,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(1),
			err:     nil,
		},
		{
			desc:    "retrieve payment ratio for non existing village",
			village: "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			payed:   uint64(0),
			pending: uint64(0),
			err:     errors.E(op, "no data found for this village", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindVillageRatio(ctx, tc.village, tc.year, tc.month)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%v' got '%v'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected payed '%d' got '%d'\n", tc.desc, tc.payed, pending))
	}
}

func TestListSectorRatios(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	var sector, cell, village = "sector", "cell", "village"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      cell,
		Sector:    sector,
		Village:   village,
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    fmt.Sprintf("%d.%s", i, cell),
				Village: fmt.Sprintf("%d.%s", i, village),
			},
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		saveProperty(t, db, p)
	}

	cases := []struct {
		desc   string
		sector string
		month  uint
		year   uint
		size   uint64
		err    error
	}{
		{
			desc:   "retrieve cells payment ratio for existing sector",
			sector: sector,
			month:  uint(time.Now().Month()),
			year:   uint(time.Now().Year()),
			size:   uint64(10),
		},
	}

	for _, tc := range cases {
		charts, err := repo.ListSectorRatios(context.TODO(), tc.sector, tc.year, tc.month)
		size := uint64(len(charts))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", tc.desc, err))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))

	}
}

func TestListCellRatios(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	var sector, cell, village = "sector", "cell", "village"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      cell,
		Sector:    sector,
		Village:   village,
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: fmt.Sprintf("%d.%s", i, village),
			},
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		saveProperty(t, db, p)
	}

	cases := []struct {
		desc  string
		cell  string
		month uint
		year  uint
		size  uint64
	}{
		{
			desc:  "retrieve villages payment ratio for existing cell",
			cell:  cell,
			month: uint(time.Now().Month()),
			year:  uint(time.Now().Year()),
			size:  uint64(10),
		},
	}

	for _, tc := range cases {
		charts, err := repo.ListCellRatios(context.TODO(), tc.cell, tc.year, tc.month)
		size := uint64(len(charts))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", tc.desc, err))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))

	}
}

func TestRetrieveSectorBalance(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindSectorBalance"

	cases := []struct {
		desc    string
		sector  string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{

		{
			desc:    "retrieve amount metrics for existing sector",
			sector:  property.Address.Sector,
			label:   property.Address.Sector,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: uint64(property.Due),
			payed:   0,
			err:     nil,
		},
		{
			desc:    "retrieve amount metrics for non-existing sector",
			sector:  "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: 0,
			payed:   0,
			err:     errors.E(op, "no data found for this sector", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindSectorBalance(ctx, tc.sector, tc.year, tc.month)
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%s' got '%s'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed'%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected pending '%d' got '%d'\n", tc.desc, tc.pending, pending))
	}
}

func TestRetrieveCellBalance(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindCellBalance"

	cases := []struct {
		desc    string
		cell    string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{

		{
			desc:    "retrieve amount metrics for existing cell",
			cell:    property.Address.Cell,
			label:   property.Address.Cell,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: uint64(property.Due),
			payed:   0,
			err:     nil,
		},
		{
			desc:    "retrieve amount metrics for non-existing cell",
			cell:    "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: 0,
			payed:   0,
			err:     errors.E(op, "no data found for this cell", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindCellBalance(ctx, tc.cell, tc.year, tc.month)
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%s' got '%s'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed'%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected pending '%d' got '%d'\n", tc.desc, tc.pending, pending))
	}
}

func TestRetrieveVillageBalance(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

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
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/statsRepository.FindVillageBalance"

	cases := []struct {
		desc    string
		village string
		month   uint
		year    uint
		label   string
		payed   uint64
		pending uint64
		err     error
	}{

		{
			desc:    "retrieve amount metrics for existing village",
			village: property.Address.Village,
			label:   property.Address.Village,
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: uint64(property.Due),
			payed:   0,
			err:     nil,
		},
		{
			desc:    "retrieve amount metrics for non-existing village",
			village: "invalid",
			month:   uint(property.CreatedAt.Month()),
			year:    uint(property.CreatedAt.Year()),
			pending: 0,
			payed:   0,
			err:     errors.E(op, "no data found for this village", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		chart, err := repo.FindVillageBalance(ctx, tc.village, tc.year, tc.month)
		payed := chart.Data["payed"]
		pending := chart.Data["pending"]
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.label, chart.Label, fmt.Sprintf("%s: expected payed '%s' got '%s'\n", tc.desc, tc.label, chart.Label))
		assert.Equal(t, tc.payed, payed, fmt.Sprintf("%s: expected payed'%d' got '%d'\n", tc.desc, tc.payed, payed))
		assert.Equal(t, tc.pending, pending, fmt.Sprintf("%s: expected pending '%d' got '%d'\n", tc.desc, tc.pending, pending))
	}
}

func TestListSectorBalances(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	var sector, cell, village = "sector", "cell", "village"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      cell,
		Sector:    sector,
		Village:   village,
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    fmt.Sprintf("%d.%s", i, cell),
				Village: fmt.Sprintf("%d.%s", i, village),
			},
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		saveProperty(t, db, p)
	}

	cases := []struct {
		desc   string
		sector string
		month  uint
		year   uint
		size   uint64
	}{
		{
			desc:   "retrieve villages payment ratio for existing cell",
			sector: sector,
			month:  uint(time.Now().Month()),
			year:   uint(time.Now().Year()),
			size:   uint64(10),
		},
	}

	for _, tc := range cases {
		charts, err := repo.ListSectorBalances(context.TODO(), tc.sector, tc.year, tc.month)
		size := uint64(len(charts))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", tc.desc, err))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))

	}
}

func TestListCellBalances(t *testing.T) {
	repo := postgres.NewStatsRepository(db)

	defer CleanDB(t, db)

	var sector, cell, village = "sector", "cell", "village"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      cell,
		Sector:    sector,
		Village:   village,
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: fmt.Sprintf("%d.%s", i, village),
			},
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		saveProperty(t, db, p)
	}

	cases := []struct {
		desc  string
		cell  string
		month uint
		year  uint
		size  uint64
	}{
		{
			desc:  "retrieve villages payment ratio for existing cell",
			cell:  cell,
			month: uint(time.Now().Month()),
			year:  uint(time.Now().Year()),
			size:  uint64(10),
		},
	}

	for _, tc := range cases {
		charts, err := repo.ListCellBalances(context.TODO(), tc.cell, tc.year, tc.month)
		size := uint64(len(charts))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", tc.desc, err))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))

	}
}
