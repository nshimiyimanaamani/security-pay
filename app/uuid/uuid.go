package uuid

import (
	"github.com/rugwirobaker/paypack-backend/app"
	uuid "github.com/satori/go.uuid"
)

var _ app.IdentityProvider = (*uuidIdentityProvider)(nil)

type uuidIdentityProvider struct{}

// New instantiates a UUID identity provider.
func New() app.IdentityProvider {
	return &uuidIdentityProvider{}
}

func (idp *uuidIdentityProvider) ID() string {
	return uuid.NewV4().String()
}
