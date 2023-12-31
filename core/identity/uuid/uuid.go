package uuid

import (
	"github.com/nshimiyimanaamani/paypack-backend/core/identity"
	uuid "github.com/satori/go.uuid"
)

var _ identity.Provider = (*uuidProvider)(nil)

type uuidProvider struct{}

// New instantiates a UUID identity provider.
func New() identity.Provider {
	return &uuidProvider{}
}

func (idp *uuidProvider) ID() string {
	return uuid.NewV4().String()
}
