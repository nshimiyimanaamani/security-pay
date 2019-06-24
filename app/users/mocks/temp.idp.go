package mocks

import (
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
)

var _ users.TempIdentityProvider = (*tempIdentityProviderMock)(nil)

type tempIdentityProviderMock struct{}

// NewTempIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewTempIdentityProvider() users.TempIdentityProvider {
	return &tempIdentityProviderMock{}
}

func (idp *tempIdentityProviderMock) TemporaryKey(id string) (string, error) {
	if id == "" {
		return "", models.ErrUnauthorizedAccess
	}

	return id, nil
}

func (idp *tempIdentityProviderMock) Identity(key string) (string, error) {
	return idp.TemporaryKey(key)
}
