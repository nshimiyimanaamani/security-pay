package jwt_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/tokens/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

var (
	secret  = "secret"
	id      = "valid"
	role    = "dev"
	account = "paypack.developers"
	invalid = "invalid"
)

func newIdentityProvider() auth.JWTProvider {
	return jwt.New(secret)
}

func TestIdentity(t *testing.T) {
	idp := newIdentityProvider()
	creds := auth.Credentials{Username: id, Role: role, Account: account}

	token, err := idp.TemporaryKey(context.Background(), creds)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "pkg/tokens/jwt.Identidy"

	cases := []struct {
		desc string
		key  string
		err  error
	}{
		{"valid key", token, nil},
		{"valid key", invalid, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)},
	}

	for _, tc := range cases {
		_, err := idp.Identity(context.Background(), tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
