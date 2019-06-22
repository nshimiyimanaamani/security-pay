package models

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	email    = "user@example.com"
	password = "password"
)

func TestUserValidate(t *testing.T){
	cases:= []struct{
		desc string
		user User
		err  error
	}{
		{"validate user with valid data", NewUser(email, password), nil},
		{"validate user with empty email", User{Email:"", Password:password}, ErrInvalidEntity},
		{"validate user with empty password", User{Email:email, Password:""}, ErrInvalidEntity},
		{"validate user with invalid email", User{Email:"userexample.com", Password:password}, ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.user.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}