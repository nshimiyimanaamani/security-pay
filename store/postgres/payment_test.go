package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestSaveTransaction(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)

	defer CleanDB(t, db)

	const op errors.Op = "store/postgres/paymentRepo.Save"

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

	defer CleanDB(t, db)

	const op errors.Op = "store/postgres/paymentRepo.RetrieveProperty"

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

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	sown := saveOwner(t, db, owner)

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
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
		id   string
		err  error
	}{
		{
			desc: "retrieve existing property",
			id:   property.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing property",
			id:   nanoid.New(nil).ID(),
			err:  errors.E(op, "property not found", errors.KindNotFound),
		},
		{
			desc: "retrieve with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "property not found", errors.KindNotFound),
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
			err:      errors.E(op, "no invoice found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.EarliestInvoice(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected error: '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
