package metrics

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Service exposes metrics usecases
type Service interface {
	FindSectorRatio(ctx context.Context, sector string) (Chart, error)
	FindCellRatio(ctx context.Context, cell string) (Chart, error)
	FindVillageRatio(ctx context.Context, village string) (Chart, error)
	ListAllSectorRatios(ctx context.Context, sector string) ([]Chart, error)
	ListAllCellRatios(ctx context.Context, cell string) ([]Chart, error)
}

// Options ...
type Options struct {
	Repo Repository
}

type service struct {
	repo Repository
}

// New metrics service
func New(opts *Options) Service {
	return &service{opts.Repo}
}

func (svc *service) FindSectorRatio(ctx context.Context, sector string) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindSectorRatio"

	chart, err := svc.repo.FindSectorRatio(ctx, sector)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}
func (svc *service) FindCellRatio(ctx context.Context, cell string) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindCellRatio"

	chart, err := svc.repo.FindCellRatio(ctx, cell)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) FindVillageRatio(ctx context.Context, village string) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindVillageRatio"

	chart, err := svc.repo.FindVillageRatio(ctx, village)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) ListAllSectorRatios(ctx context.Context, sector string) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllSectorRatios"

	charts, err := svc.repo.ListSectorRatios(ctx, sector)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}

func (svc *service) ListAllCellRatios(ctx context.Context, cell string) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllCellRatio"

	charts, err := svc.repo.ListCellRatios(ctx, cell)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}
