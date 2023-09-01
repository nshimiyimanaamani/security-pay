package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/core/nanoid"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/core/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestSavePayment(t *testing.T) {
	const op errors.Op = "store/postgres/paymentStore.Save"
	repo := postgres.NewPaymentRepository(db, nil)

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

	id := uuid.New().ID()

	cases := []struct {
		desc string
		pmt  *payment.TxRequest
		err  error
	}{
		{
			desc: "save new payment",
			pmt: &payment.TxRequest{
				ID:        id,
				Code:      property.ID,
				Amount:    invoice.Amount,
				Invoice:   invoice.ID,
				Method:    payment.MTN,
				Confirmed: false,
			},
			err: nil,
		},
		{
			desc: "save duplicate payment",
			pmt: &payment.TxRequest{
				ID:        id,
				Code:      property.ID,
				Amount:    invoice.Amount,
				Invoice:   invoice.ID,
				Method:    payment.MTN,
				Confirmed: false,
			},
			err: errors.E(op, "duplicate payment id", errors.KindAlreadyExists),
		},
		{
			desc: "save payment with invalid id",
			pmt: &payment.TxRequest{
				ID:        "invalid",
				Code:      property.ID,
				Amount:    invoice.Amount,
				Invoice:   invoice.ID,
				Method:    payment.MTN,
				Confirmed: false,
			},
			err: errors.E(op, "invalid payment entity", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Save(ctx, tc.pmt)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestFindPayment(t *testing.T) {
	const op errors.Op = "store/postgres/paymentStore.Find"

	repo := postgres.NewPaymentRepository(db, nil)
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

	invoice := invoices.Invoice{
		Amount:   property.Due,
		Property: property.ID,
		Status:   invoices.Pending,
	}

	invoice = saveInvoice(t, db, invoice)

	pmt := payment.TxRequest{
		ID:        uuid.New().ID(),
		Code:      property.ID,
		Amount:    invoice.Amount,
		Invoice:   invoice.ID,
		Method:    "mtn",
		Confirmed: false,
	}
	savePayment(t, db, pmt)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing payment",
			id:   pmt.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing payment",
			id:   nanoid.New(nil).ID(),
			err:  errors.E(op, "payment not found", errors.KindNotFound),
		},
		{
			desc: "retrieve with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "payment not found", errors.KindNotFound),
		},
	}
	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Find(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestUpdatePayment(t *testing.T) {
	const op errors.Op = "store/postgres/paymentStore.Update"

	repo := postgres.NewPaymentRepository(db, nil)
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

	invoice := invoices.Invoice{
		Amount:   property.Due,
		Property: property.ID,
		Status:   invoices.Pending,
	}

	invoice = saveInvoice(t, db, invoice)

	pmt := payment.TxRequest{
		ID:        uuid.New().ID(),
		Code:      property.ID,
		Amount:    invoice.Amount,
		Invoice:   invoice.ID,
		Method:    "mtn",
		Confirmed: false,
	}
	savePayment(t, db, pmt)

	cases := []struct {
		desc string
		pmt  payment.TxRequest
		err  error
	}{
		{
			desc: "update existing new payment",
			pmt: payment.TxRequest{
				ID:        pmt.ID,
				Confirmed: false,
			},
			err: nil,
		},
		{
			desc: "update non saved payment",
			pmt: payment.TxRequest{
				ID:        uuid.New().ID(),
				Confirmed: false,
			},
			err: errors.E(op, "payment not found", errors.KindNotFound),
		},
		{
			desc: "update payment with invalid id",
			pmt: payment.TxRequest{
				ID:        "invalid",
				Confirmed: false,
			},
			err: errors.E(op, "payment not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Update(ctx, "successuful", nil)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
