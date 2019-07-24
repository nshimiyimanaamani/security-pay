package jwt

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/stretchr/testify/assert"
)

var (
	secret  = "secret"
	id      = "valid"
	invalid = "invalid"
)

func newIdentityProvider() users.TempIdentityProvider {
	return New(secret)
}

func TestTemporaryKey(t *testing.T) {}

func TestIdentity(t *testing.T) {
	idp := newIdentityProvider()
	token, _ := idp.TemporaryKey(id)

	cases := []struct {
		desc string
		key  string
		err  error
	}{
		{"valid key", token, nil},
		{"valid key", invalid, users.ErrUnauthorizedAccess},
	}

	for _, tc := range cases {
		_, err := idp.Identity(tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
