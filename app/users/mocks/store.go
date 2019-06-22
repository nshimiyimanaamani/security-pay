package mocks

import (
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
	store "github.com/rugwirobaker/paypack-backend/store/users"
)

var _ store.UserStore = (*userStoreMock)(nil)

type userStoreMock struct {
	mu      sync.Mutex
	counter uint64
	users   map[string]models.User
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

	str.counter++
	user.ID = strconv.FormatUint(str.counter, 10)

	str.users[user.Email] = user

	return str.users[user.Email].ID, nil
}
