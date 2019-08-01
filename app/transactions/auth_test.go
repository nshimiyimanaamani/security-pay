package transactions_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/users"

	"github.com/stretchr/testify/assert"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions/mocks"
)

const (
	email = "user@emal.com"
	token = "token"
	wrong = "wrong-value"
)

func newAuthBackend(emails map[string]string) transactions.AuthBackend {
	users := mocks.NewUserService(emails)
	return transactions.NewAuthBackend(users)
}

func TestIdentify(t *testing.T) {
	auth := newAuthBackend(map[string]string{token: email})

	cases := []struct {
		desc  string
		token string
		err   error
	}{
		{
			desc:  "identify existing user",
			token: token,
			err:   nil,
		},
		{
			desc:  "identify non existant user",
			token: wrong,
			err:   users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		_, err := auth.Identity(tc.token)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
