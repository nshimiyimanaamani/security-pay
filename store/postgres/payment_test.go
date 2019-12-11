package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	props := postgres.NewPropertyStore(db)

	const op errors.Op = "postgres.paymentRepo.Save"

	defer CleanDB(t, "transactions", "properties", "owners", "users")

	var amount = 1000

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      owner,
		Due:        float64(amount),
		RecordedBy: savedUser.ID,
		Occupied:   true,
	}

	ctx := context.Background()
	_, err = props.Save(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	id := uuid.New().ID()
	method := "bk"

	cases := []struct {
		desc string
		tx   payment.Transaction
		err  error
	}{
		{
			desc: "save new transaction",
			tx:   payment.Transaction{ID: id, Code: property.ID, Amount: float64(amount), Method: method, RecordedAt: time.Now()},
			err:  nil,
		},
		{
			desc: "save duplicate transaction",
			tx:   payment.Transaction{ID: id, Code: property.ID, Amount: float64(amount), Method: method, RecordedAt: time.Now()},
			err:  errors.E(op, "duplicate transaction", errors.KindAlreadyExists),
		},
		{
			desc: "save owner with invalid id",
			tx:   payment.Transaction{ID: "invalid", Code: property.ID, Amount: float64(amount), Method: method, RecordedAt: time.Now()},
			err:  errors.E(op, "invalid transaction entity", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Save(ctx, tc.tx)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveCode(t *testing.T) {
	repo := postgres.NewPaymentRepo(db)
	props := postgres.NewPropertyStore(db)

	const op errors.Op = "postgres.paymentRepo.RetrieveCode"

	defer CleanDB(t, "transactions", "properties", "owners", "users")

	var amount = 1000

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	_, err = saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      owner,
		Due:        float64(amount),
		RecordedBy: savedUser.ID,
	}

	ctx := context.Background()
	saved, err := props.Save(ctx, property)
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
		_, err := repo.RetrieveCode(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}
