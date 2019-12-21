package randgen

import "github.com/rugwirobaker/paypack-backend/pkg/passwords"

import "context"

var _ (passwords.Generator) = (*randomGenerator)(nil)

type randomGenerator struct{}

// New creates a random password generator
func New() passwords.Generator {
	return &randomGenerator{}
}

func (gen *randomGenerator) Generate(ctx context.Context) string {
	return "password"
}
