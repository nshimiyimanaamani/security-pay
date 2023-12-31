package mocks

import (
	"context"

	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveAdmin(ctx context.Context, user users.Administrator) (users.Administrator, error) {
	const op errors.Op = "app/users/mocks/Repository.SaveAdmin"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Administrator{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) RetrieveAdmin(ctx context.Context, id string) (users.Administrator, error) {
	const op errors.Op = "app/users/mocks/Repository.RetrieveAdmin"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Administrator{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) ListAdmins(ctx context.Context, offset, limit uint64) (users.AdministratorPage, error) {
	const op errors.Op = "app/users/mocks/Repository.ListAdmins"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.AdministratorPage{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) UpdateAdminCreds(ctx context.Context, user users.Administrator) error {
	const op errors.Op = "app/users/mocks/Repository.UpdateAdminCreds"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}
