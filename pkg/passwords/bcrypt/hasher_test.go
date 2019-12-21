package bcrypt

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords"
	"github.com/stretchr/testify/assert"
)

var (
	password = "password"
	wrong    = "wrong"
)

func newHasher() passwords.Hasher {
	return New()
}

func TestHasher(t *testing.T) {
	hasher := newHasher()
	hash, _ := hasher.Hash(password)

	const op errors.Op = "pkg/passwords/bcrypt/hasher.Compare"

	cases := []struct {
		desc     string
		password string
		err      error
	}{
		{
			desc:     "compare with valid password",
			password: password,
			err:      nil,
		},
		{
			desc:     "compare an invalid password",
			password: wrong,
			err:      errors.E(op, "invalid login data: wrong password", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := hasher.Compare(tc.password, hash)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
