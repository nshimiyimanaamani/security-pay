package postgres_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	email  = "user-save@example.com"
	email2 = "user2-save@example.com"
	id     = uuid.New().ID()
)

func TestUserSave(t *testing.T) {
	repo := postgres.NewUserStore(db)

	defer CleanDB(t, "users")

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "save new user",
			user: users.User{ID: id, Username: email, Password: "pass", Cell: "admin"},
			err:  nil,
		},
		{
			desc: "save user with duplicate uuid",
			user: users.User{ID: id, Username: email2, Password: "pass", Cell: "cell1"},
			err:  users.ErrConflict,
		},
		{
			desc: "save user with duplicate email",
			user: users.User{ID: uuid.New().ID(), Username: email, Password: "pass", Cell: "cell2"},
			err:  users.ErrConflict,
		},
		{
			desc: "save user with duplicate cell",
			user: users.User{ID: uuid.New().ID(), Username: email, Password: "pass", Cell: "admin"},
			err:  users.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.Save(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUserRetrieveByID(t *testing.T) {
	user := users.User{ID: uuid.New().ID(), Username: email, Password: "pass", Cell: "cell3"}

	repo := postgres.NewUserStore(db)

	defer CleanDB(t, "users")

	_, _ = repo.Save(user)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing user", email, nil},
		{"non-existing user", "unknown@example.com", users.ErrNotFound},
	}

	for _, tc := range cases {
		_, err := repo.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
