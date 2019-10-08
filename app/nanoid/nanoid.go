package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/rugwirobaker/paypack-backend/app"
)

// nanoid settings
const (
	Alphabet = "1234567890ABCDEF"
	Length   = 8
)

var _ app.IdentityProvider = (*nanoidIdentityProvider)(nil)

type nanoidIdentityProvider struct{}

// New instantiates a UUID identity provider.
func New() app.IdentityProvider {
	return &nanoidIdentityProvider{}
}

func (idp *nanoidIdentityProvider) ID() string {
	id, _ := gonanoid.Generate(Alphabet, Length)
	return id
}
