package mocks

import (
	"context"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type mockRepository struct {
	mu      sync.Mutex
	counter uint64
	users   map[string]auth.Credentials
}

// NewRepository creates a mock instance of auth.Repository.
func NewRepository(user auth.Credentials) auth.Repository {
	users := map[string]auth.Credentials{user.Username: user}
	return &mockRepository{
		users: users,
	}
}

func (repo *mockRepository) Retrieve(ctx context.Context, username string) (auth.Credentials, error) {
	const op errors.Op = ""

	repo.mu.Lock()
	defer repo.mu.Unlock()

	empty := auth.Credentials{}

	user, ok := repo.users[username]
	if !ok {
		return empty, errors.E(op, "account not found", errors.KindNotFound)
	}

	return user, nil
}
