package jwt

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/tokens"
)

const (
	issuer   string        = "paypack"
	duration time.Duration = 10 * time.Hour
)

var _ tokens.JWTProvider = (*jwtIdentityProvider)(nil)

type jwtIdentityProvider struct {
	secret string
}

// New instantiates a JWT identity provider.
func New(secret string) tokens.JWTProvider {
	return &jwtIdentityProvider{secret}
}

// type claims struct {
// 	admin bool
// 	jwt.StandardClaims
// }

func (idp *jwtIdentityProvider) TemporaryKey(ctx context.Context, id string) (string, error) {
	const op errors.Op = "pkg/tokens/jwt.TemporaryKey"

	now := time.Now().UTC()
	exp := now.Add(duration)

	claims := jwt.StandardClaims{
		Subject:   id,
		Issuer:    issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: exp.Unix(),
	}
	return idp.jwt(claims)
}

func (idp *jwtIdentityProvider) Identity(ctx context.Context, key string) (string, error) {
	const op errors.Op = "pkg/tokens/jwt.Identidy"

	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
		}

		return []byte(idp.secret), nil
	})

	if err != nil {
		return "", errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}
	return "", nil
}

func (idp *jwtIdentityProvider) jwt(claims jwt.StandardClaims) (string, error) {
	const op errors.Op = "pkg/tokens/jwt.TemporaryKey"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(idp.secret))
	if err != nil {
		return "", errors.E(op, err)
	}
	return tokenString, nil
}
