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
	sector := "sector"
	village := "village"
	phone := "0784657205"

	cases := []struct {
		desc string
		user User
		err  error
	}{
		{"validate user with valid data(email)", User{Username: email, Password: password, Cell: cell, Sector: sector, Village: village}, nil},
		{"validate user with valid data(phone)", User{Username: phone, Password: password, Cell: cell, Sector: sector, Village: village}, nil},
		{"validate user with empty email", User{Username: "", Password: password, Cell: cell, Sector: sector, Village: village}, ErrInvalidEntity},
		{"validate user with empty password", User{Username: email, Password: "", Cell: cell, Sector: sector, Village: village}, ErrInvalidEntity},
		{"validate user with empty cell", User{Username: email, Password: password, Sector: sector, Village: village}, ErrInvalidEntity},
		{"validate user with empty sector", User{Username: email, Password: password, Cell: cell, Village: village}, ErrInvalidEntity},
		{"validate user with empty village", User{Username: email, Password: password, Cell: cell, Sector: sector}, ErrInvalidEntity},
		{"validate user with invalid email", User{Username: "userexample.com", Password: password, Cell: cell}, ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.user.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
