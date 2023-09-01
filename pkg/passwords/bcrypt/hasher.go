// Package bcrypt provides a hasher implementation utilising bcrypt.
package bcrypt

import (
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/passwords"
	"golang.org/x/crypto/bcrypt"
)

const cost int = 10

var _ passwords.Hasher = (*bcryptHasher)(nil)

type bcryptHasher struct{}

// New instantiates a bcrypt-based hasher implementation.
func New() passwords.Hasher {
	return &bcryptHasher{}
}

func (bh *bcryptHasher) Hash(pwd string) (string, error) {
	const op errors.Op = "pkg/passwords/bcrypt/hasher.Hash"

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return "", errors.E(op, err, "unexpected error", errors.KindUnexpected)
	}

	return string(hash), nil
}

func (bh *bcryptHasher) Compare(plain, hashed string) error {
	const op errors.Op = "pkg/passwords/bcrypt/hasher.Compare"

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return errors.E(op, err, "invalid login data: wrong password", errors.KindBadRequest)
		default:
			return errors.E(op, err, "unexpected error", errors.KindUnexpected)
		}
	}
	return nil
}
