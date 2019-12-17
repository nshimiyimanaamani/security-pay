package bcrypt

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	password = "password"
	wrong    = "wrong"
)

func newHasher() hasher.Hasher {
	return New()
}

func TestHasher(t *testing.T) {
	hasher := newHasher()
	hash, _ := hasher.Hash(password)

	cases := []struct {
		desc     string
		password string
		err      error
	}{
		{"compare with valid password", password, nil},
		{"compare an invalid password", wrong, bcrypt.ErrMismatchedHashAndPassword},
	}

	for _, tc := range cases {
		err := hasher.Compare(tc.password, hash)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
