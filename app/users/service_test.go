package users_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/mocks"
	"github.com/stretchr/testify/assert"
)

const wrong string = "wrong-value"

var user = users.User{Email: "user@gmail.com", Password: "password", Cell: "admin"}

func newService() users.Service {
	hasher := mocks.NewHasher()
	tempIdp := mocks.NewTempIdentityProvider()
	idp := mocks.NewIdentityProvider()
	store := mocks.NewUserStore()

	return users.New(hasher, tempIdp, idp, store)
}

func TestRegister(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{"register new user", user, nil},
		{"register existing user", user, users.ErrConflict},
		{"register new user with empty password", users.User{Email: user.Email, Cell: "admin", Password: ""}, users.ErrInvalidEntity},
		{"register new user with empty cell", users.User{Email: "new@gmail.com", Password: "new password", Cell: ""}, users.ErrInvalidEntity},
	}

	for _, tc := range cases {
		_, err := svc.Register(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestLogin(t *testing.T) {
	svc := newService()
	svc.Register(user)

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{"login with good credentials", user, nil},
		{"login with wrong e-mail", users.User{Email: wrong, Password: user.Password}, users.ErrUnauthorizedAccess},
		{"login with wrong password", users.User{Email: user.Email, Password: wrong}, users.ErrUnauthorizedAccess},
	}

	for _, tc := range cases {
		_, err := svc.Login(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestIdentify(t *testing.T) {
	svc := newService()
	svc.Register(user)
	token, _ := svc.Login(user)

	cases := []struct {
		desc  string
		token string
		err   error
	}{
		{"valid token's identity", token, nil},
		{"invalid token's identity", "", users.ErrUnauthorizedAccess},
	}
	for _, tc := range cases {
		_, err := svc.Identify(tc.token)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
