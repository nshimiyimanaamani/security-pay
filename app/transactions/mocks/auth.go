package mocks

import (
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

var _ (transactions.AuthBackend) = (*authorizationMock)(nil)

type authorizationMock struct {
	users map[string]string
}

// NewAuthBackend creates a new properties.Auth mirror
func NewAuthBackend(users map[string]string) transactions.AuthBackend {
	return &authorizationMock{users}
}

func (auth *authorizationMock) Identity(token string) (string, error) {
	if id, ok := auth.users[token]; ok {
		return id, nil
	}
	return "", users.ErrUnauthorizedAccess
}
