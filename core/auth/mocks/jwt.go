package mocks

import (
	"context"
	"fmt"
	"strings"

	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ auth.JWTProvider = (*jwtProviderMock)(nil)

type jwtProviderMock struct{}

// NewJWTProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewJWTProvider() auth.JWTProvider {
	return &jwtProviderMock{}
}

func (idp *jwtProviderMock) TemporaryKey(ctx context.Context, creds auth.Credentials) (string, error) {
	const op errors.Op = "pkg/jwt/jwtProvider.Identify"
	if creds.Username == "" {
		return "", errors.E(op, "access denied: invalid credentials", errors.KindAccessDenied)
	}
	return fmt.Sprintf("%s.%s.%s", creds.Username, creds.Role, creds.Role), nil
}

func (idp *jwtProviderMock) Identity(ctx context.Context, token string) (auth.Credentials, error) {
	const op errors.Op = "pkg/jwt/jwtProvider.Identify"

	keys := strings.Split(token, ".")

	if len(keys) < 3 {
		return auth.Credentials{}, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
	}
	creds := auth.Credentials{
		Username: keys[0],
		Account:  keys[1],
		Role:     keys[2],
	}
	if creds.Username == "" || creds.Account == "" || creds.Role == "" {
		return auth.Credentials{}, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
	}
	return creds, nil
}
