package mocks

import (
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/users"
)

var _ (users.Service) = (*userServiceMock)(nil)

type userServiceMock struct {
	mu    sync.Mutex
	users map[string]string
}

// NewUserService creates a mirror user.Service for testing
func NewUserService(users map[string]string) users.Service {
	return &userServiceMock{
		users: users,
	}
}

func (svc *userServiceMock) Register(user users.User) (string, error) {
	return "", nil
}

func (svc *userServiceMock) Login(user users.User) (string, error) {
	return "", nil
}

func (svc *userServiceMock) Identify(token string) (string, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	val, ok := svc.users[token]
	if !ok {
		return "", users.ErrUnauthorizedAccess
	}

	return val, nil
}
