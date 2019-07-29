package users

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	email := "user@example.com"
	password := "password"
	cell := "cell"

	cases := []struct {
		desc string
		user User
		err  error
	}{
		{"validate user with valid data", User{Email: email, Password: password, Cell: cell}, nil},
		{"validate user with empty email", User{Email: "", Password: password, Cell: cell}, ErrInvalidEntity},
		{"validate user with empty password", User{Email: email, Password: "", Cell: cell}, ErrInvalidEntity},
		{"validate user with invalid email", User{Email: "userexample.com", Password: password, Cell: cell}, ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.user.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestCheckCell(t *testing.T) {
	email := "user@example.com"
	password := "password"
	cell := "cell"

	cases := []struct {
		desc      string
		user      User
		populated bool
	}{
		{"check user with valid data", User{Email: email, Password: password, Cell: cell}, true},
		{"check user with invalid data", User{Email: email, Password: password, Cell: ""}, false},
	}

	for _, tc := range cases {
		populated := tc.user.CheckCell()
		assert.Equal(t, tc.populated, populated, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.populated, populated))
	}
}
