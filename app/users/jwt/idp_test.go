package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
)
var (
	secret  = "secret"
	id 	    = "valid"
	invalid = "invalid"
)

func newIdentityProvider()users.TempIdentityProvider{
	return New(secret)
}

func TestTemporaryKey(t *testing.T){}

func TestIdentity(t *testing.T){
	idp:= newIdentityProvider()
	token, _:=idp.TemporaryKey(id)

	cases:=[]struct{
		desc string
		key string
		err error
	}{
		{"valid key",token, nil},
		{"valid key", invalid, models.ErrUnauthorizedAccess},
	}

	for _, tc:= range cases{
		_, err:=idp.Identity(tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}