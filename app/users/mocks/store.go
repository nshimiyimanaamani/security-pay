package mocks

import (
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/users"
)

var _ users.Store = (*userStoreMock)(nil)

type userStoreMock struct {
	mu      sync.Mutex
	counter uint64
	users   map[string]users.User
}

//NewUserStore creates "mirror" user store
func NewUserStore() users.Store {
	return &userStoreMock{
		users: make(map[string]users.User),
	}
}

func (str *userStoreMock) RetrieveByID(email string) (users.User, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.users[email]
	if !ok {
		return users.User{}, users.ErrNotFound
	}

	return val, nil
}

func (str *userStoreMock) Save(user users.User) (string, error) {
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
