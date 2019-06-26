package postgres_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

var (
	email  = "user-save@example.com"
	email2 = "user2-save@example.com"
	id     = uuid.New().ID()
)
func TestUserSave(t *testing.T){
	
	repo := postgres.NewUserStore(db)

	cases := []struct {
		desc string
		user models.User
		err  error
	}{
		{"save new user",  models.User{ID: id,Email:email, Password:"pass",},nil},
		{"duplicate uuid", models.User{ID: id, Email:email2, Password:"pass"}, models.ErrConflict},
		{"duplicate user", models.User{ID: uuid.New().ID(), Email:email, Password:"pass"}, models.ErrConflict},
	}
	
	for _, tc:= range cases{
		_,err := repo.Save(tc.user)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUserRetrieveByID(t *testing.T){
	user:= models.User{ID:"1", Email:email, Password:"pass"}
	repo := postgres.NewUserStore(db)

	_,_=repo.Save(user)

	cases := []struct {
		desc string
		id string
		err  error
	}{
		{"retrieve existing user", email, nil},
		{"non-existing user", "unknown@example.com", models.ErrNotFound},
	}

	for _, tc:= range cases{
		_,err := repo.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}