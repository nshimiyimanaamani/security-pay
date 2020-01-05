package users

import (
	"context"
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (svc *service) RegisterAgent(ctx context.Context, user Agent) (Agent, error) {
	const op errors.Op = "app/users/service.RegisterAgent"
	if err := user.Validate(); err != nil {
		return Agent{}, errors.E(op, err)
	}

	plain := svc.pgen.Generate(ctx)

	password, err := svc.hasher.Hash(plain)
	if err != nil {
		return Agent{}, errors.E(op, err)
	}
	user.Password = password

	user.Role = Min

	now := time.Now()
	user.CreatedAt, user.UpdatedAt = now, now

	user, err = svc.repo.SaveAgent(ctx, user)
	if err != nil {
		return Agent{}, errors.E(op, err)
	}
	user.Password = plain
	return user, nil
}
func (svc *service) RetrieveAgent(ctx context.Context, id string) (Agent, error) {
	const op errors.Op = "app/users/service.RetrieveAgent"

	user, err := svc.repo.RetrieveAgent(ctx, id)
	if err != nil {
		return Agent{}, errors.E(op, err)
	}
	return user, nil
}
func (svc *service) ListAgents(ctx context.Context, offset, limit uint64) (AgentPage, error) {
	const op errors.Op = "app/users/service.ListAgents"

	page, err := svc.repo.ListAgents(ctx, offset, limit)
	if err != nil {
		return AgentPage{}, errors.E(op, err)
	}
	return page, nil
}

func (svc *service) UpdateAgentCreds(ctx context.Context, user Agent) error {
	const op errors.Op = "app/users/service.UpdateAgentCreds"

	if user.Password == "" {
		return errors.E(op, "invalid user: missing password", errors.KindBadRequest)
	}

	password, err := svc.hasher.Hash(user.Password)
	if err != nil {
		return errors.E(op, err)
	}
	user.Password = password

	user.UpdatedAt = time.Now()

	if err := svc.repo.UpdateAgentCreds(ctx, user); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) UpdateAgent(ctx context.Context, user Agent) error {
	const op errors.Op = "app/users/service.UpdateAgent"

	user.UpdatedAt = time.Now()

	if err := svc.repo.UpdateAgentDetails(ctx, user); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) DeleteAgent(ctx context.Context, id string) error {
	const op errors.Op = "app/users/service.DeleteAgent"

	err := svc.repo.DeleteAgent(ctx, id)
	if err != nil {
		errors.E(op, err)
	}

	return errors.E(op, errors.KindNotImplemented)
}
