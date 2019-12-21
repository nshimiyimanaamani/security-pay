package users

import (
	"context"
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (svc *service) RegisterAdmin(ctx context.Context, user Administrator) (Administrator, error) {
	const op errors.Op = "app/users/service.RegisterAdmin"

	if err := user.Validate(); err != nil {
		return Administrator{}, errors.E(op, err)
	}

	user.Role = Admin

	now := time.Now()
	user.CreatedAt, user.UpdatedAt = now, now

	user.Password = svc.pgen.Generate(ctx)

	return svc.repo.SaveAdmin(ctx, user)
}
func (svc *service) RetrieveAdmin(ctx context.Context, id string) (Administrator, error) {
	const op errors.Op = "app/users/service.RetrieveAdmin"

	user, err := svc.repo.RetrieveAdmin(ctx, id)
	if err != nil {
		return Administrator{}, errors.E(op, err)
	}
	return user, nil
}
func (svc *service) ListAdmins(ctx context.Context, offset, limit uint64) (AdministratorPage, error) {
	const op errors.Op = "app/users/service.ListAdmins"

	page, err := svc.repo.ListAdmins(ctx, offset, limit)
	if err != nil {
		return AdministratorPage{}, errors.E(op, err)
	}
	return page, nil
}
func (svc *service) UpdateAdminCreds(ctx context.Context, user Administrator) error {
	const op errors.Op = "app/users/service.RegisterAdmin"

	if user.Password == "" {
		return errors.E(op, "invalid user: missing password", errors.KindBadRequest)
	}

	user.UpdatedAt = time.Now()

	if err := svc.repo.UpdateAdminCreds(ctx, user); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) DeleteAdmin(ctx context.Context, id string) error {
	const op errors.Op = "app/users/service.DeleteAdmin"

	return errors.E(op, errors.KindNotImplemented)
}
