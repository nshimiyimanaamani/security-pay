package mocks

import (
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords"
)

var _ passwords.Hasher = (*hasherMock)(nil)

type hasherMock struct{}

// NewHasher creates "no-op" hasher for test purposes. This implementation will
// return secrets without changing them.
func NewHasher() passwords.Hasher {
	return &hasherMock{}
}

func (hm *hasherMock) Hash(pwd string) (string, error) {
	const op errors.Op = "pkg/passwords/hasher.Hash"
	if pwd == "" {
		return "", errors.E(op, "invalid password", errors.KindBadRequest)
	}
	return pwd, nil
}

func (hm *hasherMock) Compare(plain, hashed string) error {
	const op errors.Op = "pkg/passwords/hasher.Compare"
	if plain != hashed {
		return errors.E(op, "access denied:invalid password", errors.KindAccessDenied)
	}

	return nil
}
