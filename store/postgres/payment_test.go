package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveTransaction(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "transactions", "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner,
		Due:   float64(1000),
	}

	ctx := context.Background()
	_, err = props.Save(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	id := uuid.New().ID()
	method := "bk"

	cases := []struct {
		desc        string
		transaction payment.Transaction
		err         error
	}{
		{
			desc: "save new transaction",
			transaction: payment.Transaction{
				ID:         id,
				Code:       property.ID,
				Amount:     amount,
				Method:     method,
				RecordedAt: time.Now(),
			},
			err: nil,
		},
		{
			desc: "save duplicate transaction",
			transaction: payment.Transaction{
				ID:         id,
				Code:       property.ID,
				Amount:     amount,
				Method:     method,
				RecordedAt: time.Now(),
			},
			err: transactions.ErrConflict,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Save(ctx, tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
