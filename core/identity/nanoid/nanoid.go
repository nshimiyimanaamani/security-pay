package nanoid

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/nshimiyimanaamani/paypack-backend/core/identity"
)

// nanoid settings
const (
	Alphabet = "1234567890"
	Length   = 15
)

var _ identity.Provider = (*nanoidProvider)(nil)

type nanoidProvider struct {
	length   int
	alphabet string
}

// Config contains nanoid initialzation config
type Config struct {
	Length   int
	Alphabet string
}

// New instantiates a UUID identity provider.
func New(c *Config) identity.Provider {
	idp := &nanoidProvider{Length, Alphabet}
	if c != nil {
		idp.alphabet = c.Alphabet
		idp.length = c.Length
	}
	return idp
}

func (idp *nanoidProvider) ID() string {
	id, _ := gonanoid.Generate(idp.alphabet, idp.length)
	return id
}
