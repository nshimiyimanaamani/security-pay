package jwt

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/tokens"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	secret  = "secret"
	id      = "valid"
	invalid = "invalid"
)

func newIdentityProvider() tokens.JWTProvider {
	return New(secret)
}

func TestTemporaryKey(t *testing.T) {}

func TestIdentity(t *testing.T) {
	idp := newIdentityProvider()
	token, _ := idp.TemporaryKey(context.Background(), id)

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
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
