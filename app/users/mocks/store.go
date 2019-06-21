package mocks

import (
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
	store "github.com/rugwirobaker/paypack-backend/store/users"
	"sync"
)

var _ store.UserStore = (*userStoreMock)(nil)

type userStoreMock struct {
	mu    sync.Mutex
	users map[string]models.User
}

//NewUserStore creates "mirror" user store
func NewUserStore() store.UserStore {
	return &userStoreMock{
		users: make(map[string]models.User),
	}
}

func (str *userStoreMock) RetrieveByID(email string) (models.User, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.users[email]
	if !ok {
		return models.User{}, users.ErrNotFound
	}

	return val, nil
}

func (str *userStoreMock) Save(user models.User) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.users[user.Email]; ok {
		return "", users.ErrConflict
	}

	user.ID = "id"
	str.users[user.Email] = user

	return str.users[user.Email].ID, nil
}
