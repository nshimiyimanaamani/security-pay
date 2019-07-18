package transactions_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/transactions"

	"github.com/stretchr/testify/assert"
)

func TestTransactionValidate(t *testing.T) {
	id := "1000-4433-34343"
	amount := "5000.00"
	method := "BK"
	property := "1000-4433-34343"

	cases := []struct {
		desc  string
		tranx transactions.Transaction
		err   error
	}{
		{
			desc:  "validate user with valid data",
			tranx: transactions.Transaction{ID: id, Amount: amount, Method: method, Property: property},
			err:   nil,
		},
		{
			desc:  "validate user with invalid data",
			tranx: transactions.Transaction{},
			err:   transactions.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		err := tc.tranx.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
