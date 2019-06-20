package mocks

import "github.com/rugwirobaker/paypack-backend/app/users"

var _ users.IdentityProvider = (*identityProviderMock)(nil)

type identityProviderMock struct{}

// NewIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewIdentityProvider() users.IdentityProvider {
	return &identityProviderMock{}
}

func (idp *identityProviderMock) TemporaryKey(id string) (string, error) {
	if id == "" {
		return "", users.ErrUnauthorizedAccess
	}

	return id, nil
}

func (idp *identityProviderMock) Identity(key string) (string, error) {
	return idp.TemporaryKey(key)
}
