package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveAgent(ctx context.Context, user users.Agent) (users.Agent, error) {
	const op errors.Op = "users/Repository.SaveAgent"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Agent{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) RetrieveAgent(ctx context.Context, id string) (users.Agent, error) {
	const op errors.Op = "users/Repository.RetrieveAgent"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.Agent{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) ListAgents(ctx context.Context, offset, limit uint64) (users.AgentPage, error) {
	const op errors.Op = "users/Repository.ListAgents"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return users.AgentPage{}, errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) UpdateAgentCreds(ctx context.Context, user users.Agent) error {
	const op errors.Op = "users/Repository.UpdateAgentCreds"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}

func (repo *userRepository) UpdateAgentDetails(ctx context.Context, user users.Agent) error {
	const op errors.Op = "users/Repository.UpdateAgentDetails"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	return errors.E(op, errors.KindNotImplemented)
}
