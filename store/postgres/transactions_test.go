package postgres_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/properties"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	amount     = 2000.00
	wrongValue = "wrong"
)

func TestSave(t *testing.T) {
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	id := uuid.New().ID()
	method := "bk"

	cases := []struct {
		desc        string
		transaction transactions.Transaction
		err         error
	}{
		{
			"save new transaction",
			transactions.Transaction{
				ID:           id,
				MadeBy:       owner.ID,
				MadeFor:      property.ID,
				Amount:       amount,
				Method:       method,
				DateRecorded: time.Now(),
			},
			nil,
		},
		{
			"save duplicate transaction",
			transactions.Transaction{
				ID:           id,
				MadeBy:       owner.ID,
				MadeFor:      property.ID,
				Amount:       amount,
				Method:       method,
				DateRecorded: time.Now(),
			},
			transactions.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.Save(tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSingleTransactionRetrieveByID(t *testing.T) {
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "kcb"

	transaction := transactions.Transaction{
		ID:           uuid.New().ID(),
		MadeBy:       owner.ID,
		MadeFor:      property.ID,
		Amount:       amount,
		Method:       method,
		DateRecorded: time.Now(),
	}

	_, err = repo.Save(transaction)
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
		_, err := repo.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAll(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "mtn"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		trx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       amount,
			Method:       method,
			DateRecorded: time.Now(),
		}

		_, err := repo.Save(trx)
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
		page, err := repo.RetrieveAll(tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByProperty(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "airtel"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		trx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       amount,
			Method:       method,
			DateRecorded: time.Now(),
		}

		_, err := repo.Save(trx)
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
		page, err := repo.RetrieveByProperty(tc.property, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByMethod(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)
	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "equity"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		trx := transactions.Transaction{
			ID:           idp.ID(),
			MadeBy:       owner.ID,
			MadeFor:      property.ID,
			Amount:       amount,
			Method:       method,
			DateRecorded: time.Now(),
		}
		_, err := repo.Save(trx)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := map[string]struct {
		method string
		offset uint64
		limit  uint64
		size   uint64
	}{
		"retrieve all transactions with existing property": {
			method: method,
			offset: 0,
			limit:  n,
			size:   n,
		},
		"retrieve subset of transactions with existing property": {
			method: method,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
		},
		"retrieve transactions with non-existing property": {
			method: wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
		},
	}

	for desc, tc := range cases {
		page, err := repo.RetrieveByMethod(tc.method, tc.offset, tc.limit)
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestUpdateTransaction(t *testing.T) {
	repo := postgres.NewTransactionStore(db)
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Due:   float64(1000),
	}
	_, err = props.Save(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	method := "kcb"

	transaction := transactions.Transaction{
		ID:           uuid.New().ID(),
		MadeBy:       owner.ID,
		MadeFor:      property.ID,
		Amount:       amount,
		Method:       method,
		DateRecorded: time.Now(),
	}

	unsaved := transactions.Transaction{
		ID:           uuid.New().ID(),
		MadeBy:       owner.ID,
		MadeFor:      property.ID,
		Amount:       amount,
		Method:       method,
		DateRecorded: time.Now(),
	}

	invalid := transactions.Transaction{
		ID:           "invalid",
		MadeBy:       owner.ID,
		MadeFor:      property.ID,
		Amount:       amount,
		Method:       method,
		DateRecorded: time.Now(),
	}

	_, err = repo.Save(transaction)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc        string
		transaction transactions.Transaction
		err         error
	}{
		{
			desc:        "update existing transaction",
			transaction: transaction,
			err:         nil,
		},
		{
			desc:        "update non-existing transaction",
			transaction: unsaved,
			err:         transactions.ErrNotFound,
		},
		{
			desc:        "update with invalid data",
			transaction: invalid,
			err:         transactions.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := repo.UpdateTransaction(tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveByMonth(t *testing.T) {}

func TestRetrieveByYear(t *testing.T) {}
