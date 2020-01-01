package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveTransaction(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)

	defer CleanDB(t)

	const op errors.Op = "store/postgres/paymentRepo.Save"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))

	invoice := invoices.Invoice{Amount: property.Due, Property: property.ID, Status: invoices.Pending}
	invoice, err = saveInvoice(t, db, invoice)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))

	id := uuid.New().ID()
	method := "bk"

	cases := []struct {
		desc string
		tx   payment.Transaction
		err  error
	}{
		{
			desc: "save new transaction",
			tx:   payment.Transaction{ID: id, Code: property.ID, Amount: invoice.Amount, Invoice: invoice.ID, Method: method},
			err:  nil,
		},
		{
			desc: "save duplicate transaction",
			tx:   payment.Transaction{ID: id, Code: property.ID, Amount: invoice.Amount, Invoice: invoice.ID, Method: method},
			err:  errors.E(op, "duplicate transaction", errors.KindAlreadyExists),
		},
		{
			desc: "save owner with invalid id",
			tx:   payment.Transaction{ID: "invalid", Code: property.ID, Amount: invoice.Amount, Invoice: invoice.ID, Method: method},
			err:  errors.E(op, "invalid transaction entity", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Save(ctx, tc.tx)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error: '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveProperty(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)

	defer CleanDB(t)

	const op errors.Op = "store/postgres/paymentRepo.RetrieveProperty"

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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

	saved, err := saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing property",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing property",
			id:   nanoid.New(nil).ID(),
			err:  errors.E(op, err, "property not found", errors.KindNotFound),
		},
		{
			desc: "retrieve with malformed id",
			id:   wrongValue,
			err:  errors.E(op, err, "property not found", errors.KindNotFound),
		},
	}
	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveProperty(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error: '%s' got '%s'\n", tc.desc, tc.err, err))
	}

}

func TestOldestInvoice(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account, err := saveAccount(t, db, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	invoice := invoices.Invoice{Amount: property.Due, Property: property.ID, Status: invoices.Pending}
	invoice, err = saveInvoice(t, db, invoice)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	const op errors.Op = "store/postgres/paymentRepo.OldestInvoice"

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
			err:      errors.E(op, err, "no invoice found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.OldestInvoice(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error: '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
