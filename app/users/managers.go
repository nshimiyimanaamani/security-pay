package users

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (svc *service) RegisterManager(ctx context.Context, user Manager) (Manager, error) {
	const op errors.Op = "app/users/service.RegisterManager"

	if err := user.Validate(); err != nil {
		return Manager{}, errors.E(op, err)
	}

	user.Role = Basic

	plain := svc.pgen.Generate(ctx)

	password, err := svc.hasher.Hash(plain)
	if err != nil {
		return Manager{}, errors.E(op, err)
	}
	user.Password = password

	user, err = svc.repo.SaveManager(ctx, user)
	if err != nil {
		return Manager{}, errors.E(op, err)
	}
	user.Password = plain
	return user, nil
}

func (svc *service) RetrieveManager(ctx context.Context, id string) (Manager, error) {
	const op errors.Op = "app/users/service.RetrieveManager"

	user, err := svc.repo.RetrieveManager(ctx, id)
	if err != nil {
		return Manager{}, errors.E(op, err)
	}
	return user, nil
}
func (svc *service) ListManagers(ctx context.Context, offset, limit uint64) (ManagerPage, error) {
	const op errors.Op = "app/users/service.ListManagers"

	page, err := svc.repo.ListManagers(ctx, offset, limit)
	if err != nil {
		return ManagerPage{}, errors.E(op, err)
	}
	return page, nil
}
func (svc *service) UpdateManagerCreds(ctx context.Context, user Manager) error {
	const op errors.Op = "app/users/service.UpdateManagerCreds"

	if user.Password == "" {
		return errors.E(op, "invalid user: missing password", errors.KindBadRequest)
	}

	password, err := svc.hasher.Hash(user.Password)
	if err != nil {
		return errors.E(op, err)
	}
	user.Password = password

	if err := svc.repo.UpdateManagerCreds(ctx, user); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) DeleteManager(ctx context.Context, id string) error {
	const op errors.Op = "app/users/service.DeleteManager"

	err := svc.repo.DeleteManager(ctx, id)
	if err != nil {
		errors.E(op, err)
	}

	return errors.E(op, errors.KindNotImplemented)
}
