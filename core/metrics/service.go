package metrics

import (
	"context"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// Service exposes metrics usecases
type Service interface {
	RatioService
	BalanceService
}

// RatioService exposes payment count ratio
type RatioService interface {
	// FindSectorRatio ...
	FindSectorRatio(ctx context.Context, sector string, y, m uint) (Chart, error)

	// FindCellRatio ...
	FindCellRatio(ctx context.Context, cell string, y, m uint) (Chart, error)

	// FindVillageRatio ...
	FindVillageRatio(ctx context.Context, village string, y, m uint) (Chart, error)

	// FindVillageRatio ...
	ListAllSectorRatios(ctx context.Context, sector string, y, m uint) ([]Chart, error)

	// ListAllCellRatios ...
	ListAllCellRatios(ctx context.Context, cell string, y, m uint) ([]Chart, error)
}

// BalanceService exposes balance ratios
type BalanceService interface {
	// FindSectorBalance ...
	FindSectorBalance(ctx context.Context, sector string, y, m uint) (Chart, error)

	// FindCellBalance ...
	FindCellBalance(ctx context.Context, cell string, y, m uint) (Chart, error)

	// FindVillageBalance ...
	FindVillageBalance(ctx context.Context, village string, y, m uint) (Chart, error)

	// FindVillageBalance ...
	ListAllSectorBalances(ctx context.Context, sector string, y, m uint) ([]Chart, error)

	// ListAllCellBalances ...
	ListAllCellBalances(ctx context.Context, cell string, y, m uint) ([]Chart, error)
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

func (svc *service) FindSectorRatio(ctx context.Context, sector string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindSectorRatio"

	chart, err := svc.repo.FindSectorRatio(ctx, sector, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}
func (svc *service) FindCellRatio(ctx context.Context, cell string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindCellRatio"

	chart, err := svc.repo.FindCellRatio(ctx, cell, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) FindVillageRatio(ctx context.Context, village string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindVillageRatio"

	chart, err := svc.repo.FindVillageRatio(ctx, village, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) ListAllSectorRatios(ctx context.Context, sector string, y, m uint) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllSectorRatios"

	charts, err := svc.repo.ListSectorRatios(ctx, sector, y, m)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}

func (svc *service) ListAllCellRatios(ctx context.Context, cell string, y, m uint) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllCellRatio"

	charts, err := svc.repo.ListCellRatios(ctx, cell, y, m)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}

func (svc *service) FindSectorBalance(ctx context.Context, sector string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindSectorBalance"

	chart, err := svc.repo.FindSectorBalance(ctx, sector, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}
func (svc *service) FindCellBalance(ctx context.Context, cell string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindCellBalance"

	chart, err := svc.repo.FindCellBalance(ctx, cell, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) FindVillageBalance(ctx context.Context, village string, y, m uint) (Chart, error) {
	const op errors.Op = "app/metrics/service.FindVillageBalance"

	chart, err := svc.repo.FindVillageBalance(ctx, village, y, m)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) ListAllSectorBalances(ctx context.Context, sector string, y, m uint) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllSectorBalances"

	charts, err := svc.repo.ListSectorBalances(ctx, sector, y, m)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}

func (svc *service) ListAllCellBalances(ctx context.Context, cell string, y, m uint) ([]Chart, error) {
	const op errors.Op = "app/metrics/service.ListAllCellBalance"

	charts, err := svc.repo.ListCellBalances(ctx, cell, y, m)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return charts, nil
}
