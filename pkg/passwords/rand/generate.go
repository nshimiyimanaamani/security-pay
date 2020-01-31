package rand

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords"
)

const length int = 8

var _ (passwords.Generator) = (*randomGenerator)(nil)

type randomGenerator struct{}

// New creates a random password generator
func New() passwords.Generator {
	return &randomGenerator{}
}

func (g *randomGenerator) Generate(ctx context.Context) (pass string, err error) {
	const op errors.Op = "pkg/passwords/rand/Generator.Generate"

	b, err := g.GenerateRandomBytes(length)
	if err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (g *randomGenerator) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
