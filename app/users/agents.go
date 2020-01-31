package users

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (svc *service) RegisterAgent(ctx context.Context, user Agent) (Agent, error) {
	const op errors.Op = "app/users/service.RegisterAgent"
	if err := user.Validate(); err != nil {
		return Agent{}, errors.E(op, err)
	}

	plain, err := svc.pgen.Generate(ctx)

	if err != nil {
		return Agent{}, errors.E(op, err)
	}

	encrypted, err := svc.encrypter.Encrypt(plain)
	if err != nil {
		return Agent{}, errors.E(op, err)
	}

	user.Password = string(encrypted)
	user.Role = Min

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
	user.Password, err = svc.encrypter.Decrypt([]byte(user.Password))
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

	b, err := svc.encrypter.Encrypt(user.Password)
	if err != nil {
		return errors.E(op, err)
	}
	user.Password = string(b)

	if err := svc.repo.UpdateAgentCreds(ctx, user); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) UpdateAgent(ctx context.Context, user Agent) error {
	const op errors.Op = "app/users/service.UpdateAgent"

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
