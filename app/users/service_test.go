package users_test

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/mocks"
	//store "github.com/rugwirobaker/paypack-backend/store/users"
)

const wrong string = "wrong-value"

var user = models.NewUser("user@gmail.com", "password")

func newService() users.Service {
	idp := mocks.NewIdentityProvider()
	hasher:= mocks.NewHasher()
	config:= mocks.NewConfiguration()
	store:= mocks.NewUserStore()

	return users.New(idp, config, hasher, store)
}

func TestRegister(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc string
		user models.User
		err  error
	}{
		{"register new user", user, nil},
		{"register existing user", user, users.ErrConflict},
		{"register new user with empty password", models.NewUser(user.Email, ""), users.ErrInvalidEntity},
	}

	for _, tc := range cases {
		_, err:= svc.Register(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestLogin(t *testing.T){
	svc := newService()
	svc.Register(user)

	cases := []struct {
		desc string
		user models.User
		err  error
	}{
		{"login with good credentials", user, nil},
		{"login with wrong e-mail", models.NewUser(wrong, user.Password), users.ErrUnauthorizedAccess},
		{"login with wrong password", models.NewUser(user.Email, wrong), users.ErrUnauthorizedAccess },
	}

	for _, tc := range cases {
		_, err:= svc.Login(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestIdentify(t *testing.T){
	svc := newService()
	svc.Register(user)
	token, _ := svc.Login(user)

	cases := []struct {
		desc string
		token string
		err  error
	}{
		{"valid token's identity", token, nil},
		{"invalid token's identity", "", users.ErrUnauthorizedAccess },
	}
	for _, tc := range cases {
		_, err:= svc.Identify(tc.token)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}