package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestFindInvoice(t *testing.T) {
	const op errors.Op = "store/postgres/invoices.Find"

	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t, db)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	//save owner
	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

	cases := []struct {
		desc string
		id   uint64
		err  error
	}{
		{
			desc: "retrieve existing invoice record",
			id:   1,
			err:  nil,
		},
		{
			desc: "retrieve non existing invoice record",
			id:   10,
			err:  errors.E(op, "invoice not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Find(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestListAll(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t, db)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	//save owner
	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

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
		page, err := repo.All(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}

func TestListPending(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t, db)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	//save owner
	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	owner = saveOwner(t, db, owner)

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

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
		page, err := repo.Pending(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}

func TestListPayed(t *testing.T) {
	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t, db)

	// save account
	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

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
	agent = saveAgent(t, db, agent)

	//save owner
	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	owner = saveOwner(t, db, owner)

	//save property
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

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
		page, err := repo.Payed(context.Background(), tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected invoices: '%d' got '%d'\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total:'%d; got '%d'\n", tc.desc, tc.total, page.Total))

	}
}

func TestEarliest(t *testing.T) {
	const op errors.Op = "store/postgres/invoiceRepository.Earliest"

	repo := postgres.NewInvoiceRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
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

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

	invoice := invoices.Invoice{Amount: property.Due, Property: property.ID, Status: invoices.Pending}
	invoice = saveInvoice(t, db, invoice)

	cases := []struct {
		desc     string
		property string
		err      error
	}{
		{
			desc:     "retrieve oldest invoice for existing property",
			property: invoice.Property,
			err:      nil,
		},
		{
			desc:     "retrieve oldest invoice for existing property",
			property: "invalid",
			err:      errors.E(op, "no invoice found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Earliest(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error: '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
