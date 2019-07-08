package postgres_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	amount     = "2000.00"
	wrongValue = "wrong"
)

func TestSave(t *testing.T) {
	repo := postgres.NewTransactionStore(db)
	defer CleanDB(t, "transactions")

	property := uuid.New().ID()
	id := uuid.New().ID()
	method := "bk"

	cases := []struct {
		desc        string
		transaction models.Transaction
		err         error
	}{
		{
			"save new transaction",
			models.Transaction{
				ID:       id,
				Property: property,
				Amount:   amount,
				Method:   method,
			},
			nil,
		},
		{
			"save duplicate transaction",
			models.Transaction{
				ID:       id,
				Property: uuid.New().ID(),
				Amount:   "4000.00",
				Method:   "bk",
			},
			models.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.Save(tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSinglePropertyRetrieveByID(t *testing.T) {
	repo := postgres.NewTransactionStore(db)
	defer CleanDB(t, "transactions")

	property := uuid.New().ID()
	method := "kcb"

	transaction := models.Transaction{
		ID:       uuid.New().ID(),
		Property: property,
		Amount:   amount,
		Method:   method,
	}

	id, _ := repo.Save(transaction)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing transaction", id, nil},
		{"retrieve non existing transaction", uuid.New().ID(), models.ErrNotFound},
		{"retrieve with malformed id", wrongValue, models.ErrNotFound},
	}

	for _, tc := range cases {
		_, err := repo.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAll(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	defer CleanDB(t, "transactions")

	property := uuid.New().ID()
	method := "mtn"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		t := models.Transaction{
			ID:       idp.ID(),
			Property: property,
			Amount:   amount,
			Method:   method,
		}

		repo.Save(t)
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
		page := repo.RetrieveAll(tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByProperty(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	defer CleanDB(t, "transactions")

	property := uuid.New().ID()
	method := "airtel"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		t := models.Transaction{
			ID:       idp.ID(),
			Property: property,
			Amount:   amount,
			Method:   method,
		}

		repo.Save(t)
	}

	cases := map[string]struct {
		property string
		offset   uint64
		limit    uint64
		size     uint64
	}{
		"retrieve all transactions with existing property": {
			property: property,
			offset:   0,
			limit:    n,
			size:     n,
		},
		"retrieve subset of transactions with existing property": {
			property: property,
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
		page := repo.RetrieveByProperty(tc.property, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByMethod(t *testing.T) {
	idp := uuid.New()
	repo := postgres.NewTransactionStore(db)
	defer CleanDB(t, "transactions")

	property := uuid.New().ID()
	method := "equity"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		t := models.Transaction{
			ID:       idp.ID(),
			Property: property,
			Amount:   amount,
			Method:   method,
		}

		repo.Save(t)
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
		page := repo.RetrieveByMethod(tc.method, tc.offset, tc.limit)
		size := uint64(len(page.Transactions))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
	}
}

func TestRetrieveByMonth(t *testing.T) {}

func TestRetrieveByYear(t *testing.T) {}
