package users

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

func (svc *service) RegisterDeveloper(ctx context.Context, user Developer) (Developer, error) {
	const op errors.Op = "app/users/service.RegisterDeveloper"

	if err := user.Validate(); err != nil {
		return Developer{}, errors.E(op, err)
	}

	user, err := svc.repo.SaveDeveloper(ctx, user)
	if err != nil {
		return Developer{}, errors.E(op, err)
	}
	return user, nil
}

func (svc *service) RetrieveDeveloper(ctx context.Context, id string) (Developer, error) {
	const op errors.Op = "app/users/service.RetrieveDeveloper"

	user, err := svc.repo.RetrieveDeveloper(ctx, id)
	if err != nil {
		return Developer{}, errors.E(op, err)
	}
	return user, nil
}

func (svc *service) ListDevelopers(ctx context.Context, offset, limit uint64) (DeveloperPage, error) {
	const op errors.Op = "app/users/service.ListDevelopers"

	page, err := svc.repo.ListDevelopers(ctx, offset, limit)
	if err != nil {
		return DeveloperPage{}, errors.E(op, err)
	}
	return page, nil
}
func (svc *service) UpdateDeveloperCreds(ctx context.Context, user Developer) error {
	const op errors.Op = "app/users/service.UpdateDeveloperCreds"
	err := svc.repo.UpdateDeveloperCreds(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) DeleteDeveloper(ctx context.Context, id string) error {
	const op errors.Op = "app/users/service.DeleteDeveloper"

	return errors.E(op, errors.KindNotImplemented)
}
