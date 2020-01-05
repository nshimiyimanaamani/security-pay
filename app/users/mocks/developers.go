package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveDeveloper(ctx context.Context, user users.Developer) (users.Developer, error) {
	const op errors.Op = "app/users/mocks/Repository.SaveDeveloper"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Developer{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) RetrieveDeveloper(ctx context.Context, id string) (users.Developer, error) {
	const op errors.Op = "app/users/mocks/Repository.RetrieveDeveloper"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Developer{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) ListDevelopers(ctx context.Context, offset, limit uint64) (users.DeveloperPage, error) {
	const op errors.Op = "app/users/mocks/Repository.ListDevelopers"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.DeveloperPage{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) UpdateDeveloperCreds(ctx context.Context, user users.Developer) error {
	const op errors.Op = "app/users/mocks/Repository.UpdateDeveloperCreds"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) DeleteDeveloper(ctx context.Context, id string) error {
	const op errors.Op = "app/users/mocks/Repository.DeleteDeveloper"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}
