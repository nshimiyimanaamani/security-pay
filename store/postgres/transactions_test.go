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
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/tools"

	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	amount     = 2000.00
	wrongValue = "wrong"
)

func TestSingleTransactionRetrieveByID(t *testing.T) {
	repo := postgres.NewTransactionRepository(db)

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
		ID:        uuid.New().ID(),
		Fname:     "rugwiro",
		Lname:     "james",
		Phone:     "0784677882",
		Namespace: account.ID,
	}
	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: owner.ID},
		Due:        float64(1000),
		Namespace:  account.ID,
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

	invoice := retrieveInvoice(t, db, property.ID)

	method := "kcb"

	transaction := transactions.Transaction{
		ID:      uuid.New().ID(),
		OwnerID: owner.ID,
		MadeFor: property.ID,
		Amount:  invoice.Amount,
		Method:  method,
		Invoice: invoice.ID,
	}
	saveTx(t, db, transaction)

	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByID"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing transaction",
			id:   transaction.ID,
			err:  nil,
		},
		{
			desc: "retrieve non existing transaction",
			id:   uuid.New().ID(),
			err:  errors.E(op, "transaction not found", errors.KindNotFound),
		},
		{
			desc: "retrieve with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "transaction not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveByID(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestRetrieveAll(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
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
		ID:        uuid.New().ID(),
		Fname:     "rugwiro",
		Lname:     "james",
		Phone:     "0784677882",
		Namespace: account.ID,
	}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		property := properties.Property{
			ID:         nanoid.New(nil).ID(),
			Owner:      properties.Owner{ID: owner.ID},
			Due:        float64(1000),
			Namespace:  account.ID,
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}
		property = saveProperty(t, db, property)

		invoice := retrieveInvoice(t, db, property.ID)

		tx := transactions.Transaction{
			ID:      idp.ID(),
			OwnerID: owner.ID,
			MadeFor: property.ID,
			Amount:  invoice.Amount,
			Method:  "mtn",
			Invoice: invoice.ID,
		}
		saveTx(t, db, tx)
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
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByProperty(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "developers",
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
		ID:        uuid.New().ID(),
		Fname:     "rugwiro",
		Lname:     "james",
		Phone:     "0784677882",
		Namespace: account.ID,
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
		ForRent:    true,
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

	n := uint64(10)

	begin := tools.BeginningOfMonth()

	for i := uint64(1); i <= n; i++ {
		invoice := invoices.Invoice{
			Amount:    property.Due,
			Property:  property.ID,
			Status:    invoices.Pending,
			CreatedAt: tools.AddMonth(begin, int(i)),
			UpdatedAt: tools.AddMonth(begin, int(i)),
		}
		invoice = saveInvoice(t, db, invoice)

		tx := transactions.Transaction{
			ID:      idp.ID(),
			OwnerID: owner.ID,
			MadeFor: property.ID,
			Amount:  invoice.Amount,
			Invoice: invoice.ID,
			Method:  "airtel",
		}
		saveTx(t, db, tx)
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
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByPropertyR(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "developers",
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
		ID:        uuid.New().ID(),
		Fname:     "rugwiro",
		Lname:     "james",
		Phone:     "0784677882",
		Namespace: account.ID,
	}

	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: owner.ID},
		Due:        float64(1000),
		Namespace:  account.ID,
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}

	property = saveProperty(t, db, property)

	begin := tools.BeginningOfMonth()

	n := uint64(10)

	for i := uint64(1); i <= n; i++ {
		invoice := invoices.Invoice{
			Amount:    property.Due,
			Property:  property.ID,
			Status:    invoices.Pending,
			CreatedAt: tools.AddMonth(begin, int(i)),
			UpdatedAt: tools.AddMonth(begin, int(i)),
		}
		invoice = saveInvoice(t, db, invoice)

		tx := transactions.Transaction{
			ID:      idp.ID(),
			OwnerID: owner.ID,
			MadeFor: property.ID,
			Amount:  invoice.Amount,
			Invoice: invoice.ID,
			Method:  "airtel",
		}
		saveTx(t, db, tx)
	}

	cases := map[string]struct {
		property string
		size     uint64
	}{
		"retrieve all transactions with existing property": {
			property: property.ID,
			size:     n,
		},
		"retrieve transactions with non-existing property": {
			property: uuid.New().ID(),
			size:     0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveByPropertyR(ctx, tc.property)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByMethod(t *testing.T) {
	t.Skip()

	idp := uuid.New()
	repo := postgres.NewTransactionRepository(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
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
		ID:        uuid.New().ID(),
		Fname:     "rugwiro",
		Lname:     "james",
		Phone:     "0784677882",
		Namespace: account.ID,
	}
	owner = saveOwner(t, db, owner)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {

		property := properties.Property{
			ID:         nanoid.New(nil).ID(),
			Owner:      properties.Owner{ID: owner.ID},
			Due:        float64(1000),
			Namespace:  account.ID,
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}
		property = saveProperty(t, db, property)

		invoice := retrieveInvoice(t, db, property.ID)

		tx := transactions.Transaction{
			ID:      idp.ID(),
			OwnerID: owner.ID,
			MadeFor: property.ID,
			Amount:  invoice.Amount,
			Invoice: invoice.ID,
			Method:  "mtn",
		}

		saveTx(t, db, tx)
	}

	cases := map[string]struct {
		method string
		offset uint64
		limit  uint64
		size   uint64
	}{}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveByMethod(ctx, tc.method, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}
