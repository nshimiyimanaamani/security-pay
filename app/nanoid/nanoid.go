package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/rugwirobaker/paypack-backend/app"
)

// nanoid settings
const (
	Alphabet = "1234567890"
	Length   = 15
)

var _ app.IdentityProvider = (*nanoidIdentityProvider)(nil)

type nanoidIdentityProvider struct {
	length   int
	alphabet string
}

// Config contains nanoid initialzation config
type Config struct {
	Length   int
	Alphabet string
}

// New instantiates a UUID identity provider.
func New(c *Config) app.IdentityProvider {
	idp := &nanoidIdentityProvider{Length, Alphabet}
	if c != nil {
		idp.alphabet = c.Alphabet
		idp.length = c.Length
	}
	return idp
}

func (idp *nanoidIdentityProvider) ID() string {
	id, _ := gonanoid.Generate(idp.alphabet, idp.length)
	return id
}
