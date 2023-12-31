package accounts_test

import (
	"fmt"
	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateAccount(t *testing.T) {
	const op errors.Op = "app/accounts/account.Validate"

	cases := []struct {
		desc    string
		account accounts.Account
		err     error
	}{
		{
			desc:    "validate valid account",
			account: accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     nil,
		},
		{
			desc:    "validate account with missing name",
			account: accounts.Account{ID: "paypack.developers", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing name", errors.KindBadRequest),
		},
		{
			desc:    "validate account with missing type",
			account: accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10},
			err:     errors.E(op, "invalid account: missing type", errors.KindBadRequest),
		},
		{
			desc:    "validate account with missing account id(sector)",
			account: accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs},
			err:     errors.E(op, "invalid account: missing id", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.account.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
