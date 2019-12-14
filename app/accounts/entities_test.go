package accounts_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	const op errors.Op = "accounts/account.Validate"
	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{
		{
			desc:    "validate valid account",
			account: accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     nil,
		},
		{
			desc:    "validate account with missing name",
			account: accounts.Account{NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing name", errors.KindBadRequest),
		},
		{
			desc:    "validate account with missing type",
			account: accounts.Account{Name: "remera", NumberOfSeats: 10},
			err:     errors.E(op, "invalid account: missing type", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.account.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
