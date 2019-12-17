package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ auth.JWTProvider = (*jwtProviderMock)(nil)

type jwtProviderMock struct{}

// NewTempIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewTempIdentityProvider() auth.JWTProvider {
	return &jwtProviderMock{}
}

func (idp *jwtProviderMock) TemporaryKey(ctx context.Context, id string) (string, error) {
	const op errors.Op = "pkg/jwt/jwtProvider.Identify"
	if id == "" {
		return "", errors.E(op, "access denied: invalid credentials", errors.KindAccessDenied)
	}

	return id, nil
}

func (idp *jwtProviderMock) Identity(ctx context.Context, key string) (string, error) {
	return idp.TemporaryKey(ctx, key)
}
