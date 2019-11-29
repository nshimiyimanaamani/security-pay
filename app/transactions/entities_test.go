package transactions_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/transactions"

	"github.com/stretchr/testify/assert"
)

var wrong = "wrong_value"

func TestTransactionValidate(t *testing.T) {
	id := "1000-4433-34343"
	amount := 5000.00
	method := "BK"
	property := "1000-4433-34343"
	owner := "1000-4433-34343"

	cases := []struct {
		desc string
		tx   transactions.Transaction
		err  error
	}{
		{
			desc: "validate user with valid data",
			tx: transactions.Transaction{
				ID: id, Amount: amount, Method: method, MadeBy: owner, MadeFor: property,
			},
			err: nil,
		},
		{
			desc: "validate user with invalid data",
			tx:   transactions.Transaction{},
			err:  transactions.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		err := tc.tx.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
