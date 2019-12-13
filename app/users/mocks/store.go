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

func (str *userStoreMock) Save(user users.User) (users.User, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	empty := users.User{}

	if _, ok := str.users[user.Username]; ok {
		return empty, users.ErrConflict
	}

	str.counter++
	user.ID = strconv.FormatUint(str.counter, 10)

	str.users[user.Username] = user

	// technical debt
	user.ID = user.Username

	return user, nil
}
