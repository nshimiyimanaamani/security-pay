package jwt

import (
	"time"
	"github.com/rugwirobaker/paypack-backend/app/users"
)
const (
	issuer   string        = "paypack"
	duration time.Duration = 10 * time.Hour
)

var _ users.IdentityProvider = (*jwtIdentityProvider)(nil)

type jwtIdentityProvider struct{
	secret string
}

func (idp *jwtIdentityProvider)TemporaryKey(string) (string, error){
	return "", nil
}

// Identity extracts the entity identifier given its secret key.
func (idp *jwtIdentityProvider)Identity(string) (string, error){
	return "", nil
}
