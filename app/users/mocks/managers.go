package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveManager(ctx context.Context, user users.Manager) (users.Manager, error) {
	const op errors.Op = "users/Repository.SaveManager"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Manager{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) RetrieveManager(ctx context.Context, id string) (users.Manager, error) {
	const op errors.Op = "users/Repository.RetrieveManager"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Manager{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) ListManagers(ctx context.Context, offset, limit uint64) (users.ManagerPage, error) {
	const op errors.Op = "users/Repository.ListManagers"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.ManagerPage{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) UpdateManagerCreds(ctx context.Context, user users.Manager) error {
	const op errors.Op = "users/Repository.UpdateManagerCreds"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}
