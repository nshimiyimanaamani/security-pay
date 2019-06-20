package users_test

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/mocks"
)

const wrong string = "wrong-value"

var user = models.NewOperator("operator@gmail.com", "password")

func newService() users.Service {
	idp := mocks.NewIdentityProvider()
	config:= mocks.NewConfiguration()

	return users.New(idp, config)
}

func TestRegister(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc string
		user models.Operator
		err  error
	}{
		{"register new user", user, nil},
	}

	for _, tc := range cases {
		user_id, err:= svc.Register(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.NotEqual(t, len(user_id), 0, "expected id length > 0")
	}
}

func TestLogin(t *testing.T){
	svc := newService()

	cases := []struct {
		desc string
		user models.Operator
		err  error
	}{
		{"login with good credentials", user, nil},
	}

	for _, tc := range cases {
		user_id, err:= svc.Login(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.NotEqual(t, len(user_id), 0, "expected id length > 0")
	}
}

func TestIdentify(t *testing.T){
	svc := newService()
	token, _ := svc.Login(user)

	cases := []struct {
		desc string
		token string
		err  error
	}{
		{"valid token's identity", token, nil},
	}
	for _, tc := range cases {
		_, err:= svc.Identify(tc.token)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}