package jwt

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

const (
	issuer   string        = "paypack"
	duration time.Duration = 10 * time.Hour
)

var _ auth.JWTProvider = (*jwtIdentityProvider)(nil)

type jwtIdentityProvider struct {
	secret string
}

// New instantiates a JWT identity provider.
func New(secret string) auth.JWTProvider {
	return &jwtIdentityProvider{secret}
}

// Claims are custom jwt
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Account  string `json:"account"`
	jwt.StandardClaims
}

func (idp *jwtIdentityProvider) TemporaryKey(ctx context.Context, creds auth.Credentials) (string, error) {
	const op errors.Op = "pkg/tokens/jwt.TemporaryKey"

	now := time.Now().UTC()
	exp := now.Add(duration)

	claims := &Claims{
		Username: creds.Username,
		Role:     creds.Role,
		Account:  creds.Account,
		StandardClaims: jwt.StandardClaims{
			Subject:   creds.Username,
			Issuer:    issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: exp.Unix(),
		},
	}
	return idp.jwt(claims)
}

func (idp *jwtIdentityProvider) Identity(ctx context.Context, key string) (auth.Credentials, error) {
	const op errors.Op = "pkg/tokens/jwt.Identidy"

	var claims Claims

	token, err := jwt.ParseWithClaims(key, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
		}

		return []byte(idp.secret), nil
	})

	if err != nil {
		return auth.Credentials{}, errors.E(op, err, "access denied: invalid token", errors.KindAccessDenied)
	}

	if !token.Valid {
		return auth.Credentials{}, errors.E(op, "access denied: invalid token", errors.KindAccessDenied)
	}
	creds := auth.Credentials{
		Username: claims.Username,
		Account:  claims.Account,
		Role:     claims.Role,
	}
	return creds, nil
}

func (idp *jwtIdentityProvider) jwt(claims jwt.Claims) (string, error) {
	const op errors.Op = "pkg/tokens/jwt.TemporaryKey"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(idp.secret))
	if err != nil {
		return "", errors.E(op, err)
	}
	return tokenString, nil
}
