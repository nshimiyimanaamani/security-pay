package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/users"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	amount     = 2000.00
	wrongValue = "wrong"
)

func TestSingleTransactionRetrieveByID(t *testing.T) {
	repo := postgres.NewTransactionRepository(db)

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

	savedAgent, err := saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	sown, err := saveOwner(t, db, owner)

	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: sown.ID},
		Due:        float64(1000),
		RecordedBy: savedAgent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	//save invoice
	invoice := invoices.Invoice{Amount: property.Due, Property: property.ID, Status: invoices.Pending}
	invoice, err = saveInvoice(t, db, invoice)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "kcb"

	transaction := transactions.Transaction{
		ID:           uuid.New().ID(),
		MadeBy:       owner.ID,
		MadeFor:      property.ID,
		Amount:       invoice.Amount,
		Method:       method,
		Invoice:      invoice.ID,
		DateRecorded: time.Now(),
	}

	_, err = saveTx(t, db, transaction)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing transaction", transaction.ID, nil},
		{"retrieve non existing transaction", uuid.New().ID(), transactions.ErrNotFound},
		{"retrieve with malformed id", wrongValue, transactions.ErrNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected err '%s' got '%s'\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAll(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
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
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: owner.ID},
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	//save invoice
	invoice := invoices.Invoice{Amount: property.Due, Property: property.ID, Status: invoices.Pending}
	invoice, err = saveInvoice(t, db, invoice)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "mtn"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		invoice := invoices.Invoice{
			Amount:   property.Due,
			Property: property.ID,
			Status:   invoices.Pending,
		}
		invoice, err = saveInvoice(t, db, invoice)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		tx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       invoice.Amount,
			Method:       method,
			Invoice:      invoice.ID,
			DateRecorded: time.Now(),
		}
		_, err = saveTx(t, db, tx)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
	}{
		"retrieve all transactions": {
			offset: 0,
			limit:  n,
			size:   n,
		},
		"retreive a subset of all transactions": {
			offset: 0,
			limit:  n / 2,
			size:   n / 2,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveAll(ctx, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByProperty(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
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

	savedAgent, err := saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	sown, err := saveOwner(t, db, owner)

	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: sown.ID},
		Due:        float64(1000),
		RecordedBy: savedAgent.Telephone,
		Occupied:   true,
	}

	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "airtel"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		invoice := invoices.Invoice{
			Amount:   property.Due,
			Property: property.ID,
			Status:   invoices.Pending,
		}
		invoice, err = saveInvoice(t, db, invoice)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		tx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       invoice.Amount,
			Invoice:      invoice.ID,
			Method:       method,
			DateRecorded: time.Now(),
		}

		_, err = saveTx(t, db, tx)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := map[string]struct {
		property string
		offset   uint64
		limit    uint64
		size     uint64
	}{
		"retrieve all transactions with existing property": {
			property: property.ID,
			offset:   0,
			limit:    n,
			size:     n,
		},
		"retrieve subset of transactions with existing property": {
			property: property.ID,
			offset:   n / 2,
			limit:    n,
			size:     n / 2,
		},
		"retrieve transactions with non-existing property": {
			property: uuid.New().ID(),
			offset:   0,
			limit:    n,
			size:     0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveByProperty(ctx, tc.property, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByMethod(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
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

	savedAgent, err := saveAgent(t, db, agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: sown.ID},
		Due:        float64(1000),
		RecordedBy: savedAgent.Telephone,
		Occupied:   true,
	}
	property, err = saveProperty(t, db, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "equity"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		invoice := invoices.Invoice{
			Amount:   property.Due,
			Property: property.ID,
			Status:   invoices.Pending,
		}
		invoice, err = saveInvoice(t, db, invoice)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		tx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       invoice.Amount,
			Invoice:      invoice.ID,
			Method:       method,
			DateRecorded: time.Now(),
		}

		_, err = saveTx(t, db, tx)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := map[string]struct {
		method string
		offset uint64
		limit  uint64
		size   uint64
	}{}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveByPeriod(ctx, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}
