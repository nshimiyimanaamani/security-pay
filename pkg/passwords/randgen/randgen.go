package randgen

import (
	"context"
	"math/rand"

	"github.com/rugwirobaker/paypack-backend/pkg/passwords"
)

const length int = 8

var dict = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var _ (passwords.Generator) = (*randomGenerator)(nil)

type randomGenerator struct{}

// New creates a random password generator
func New() passwords.Generator {
	return &randomGenerator{}
}

func (gen *randomGenerator) Generate(ctx context.Context) string {

	b := make([]rune, length)
	for i := range b {
		b[i] = dict[rand.Intn(len(dict))]
	}
	return string(b)
}
