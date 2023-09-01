package transactions_test

import (
	"fmt"
	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"

	"github.com/stretchr/testify/assert"
)

var wrong = "wrong_value"

func TestTransactionValidate(t *testing.T) {
	id := "1000-4433-34343"
	amount := 5000.00
	method := "BK"
	property := "1000-4433-34343"
	owner := "1000-4433-34343"

	const op errors.Op = "app/transactions/transaction.Validate"

	cases := []struct {
		desc string
		tx   transactions.Transaction
		err  error
	}{
		{
			desc: "validate user with valid data",
			tx: transactions.Transaction{
				ID: id, Amount: amount, Method: method, OwnerID: owner, MadeFor: property,
			},
			err: nil,
		},
		{
			desc: "validate user with invalid data",
			tx:   transactions.Transaction{},
			err:  errors.E(op, "invalid transaction", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.tx.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
