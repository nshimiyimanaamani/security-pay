package auth

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords"
)

var _ (Service) = (*service)(nil)

// Service aggregates Authentication usecases
type Service interface {
	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(ctx context.Context, user Credentials) (string, error)

	// Identify validates user's token. If token is valid, user's id
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error values are returned in response.
	Identify(ctx context.Context, token string) (Credentials, error)
}

// Options minimises New function signature
type Options struct {
	Hasher passwords.Hasher
	Repo   Repository
	JWT    JWTProvider
}

type service struct {
	hasher passwords.Hasher
	repo   Repository
	jwt    JWTProvider
}

// New creates a new auth.Service instance
func New(opts *Options) Service {
	return &service{
		hasher: opts.Hasher,
		repo:   opts.Repo,
		jwt:    opts.JWT,
	}
}

// encode creds including role
func (svc *service) Login(ctx context.Context, user Credentials) (string, error) {
	const op errors.Op = "app/auth/service.Login"

	creds, err := svc.repo.Retrieve(ctx, user.Username)
	if err != nil {
		return "", errors.E(op, err)
	}
	if err := svc.hasher.Compare(user.Password, creds.Password); err != nil {
		return "", errors.E(op, err)
	}

	token, err := svc.jwt.TemporaryKey(ctx, creds)
	if err != nil {
		return "", errors.E(op, err)
	}
	return token, nil
}

// must return creds
func (svc *service) Identify(ctx context.Context, token string) (Credentials, error) {
	const op errors.Op = "app/auth/service.Identify"

	creds, err := svc.jwt.Identity(ctx, token)
	if err != nil {
		return Credentials{}, errors.E(op, err)
	}
	return creds, nil
}
