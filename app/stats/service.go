package stats

import "context"

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Service exposes stats usecases
type Service interface {
	SectorPayRatio(ctx context.Context, sector string) (Chart, error)
	CellPayRatio(ctx context.Context, cell string) (Chart, error)
	VillagePayRatio(ctx context.Context, village string) (Chart, error)
}

// Options ...
type Options struct {
	Repo Repository
}

type service struct {
	repo Repository
}

// New stats service
func New(opts *Options) Service {
	return &service{opts.Repo}
}

func (svc *service) SectorPayRatio(ctx context.Context, sector string) (Chart, error) {
	const op errors.Op = "app/stats/service.SectorPayRatio"

	chart, err := svc.repo.RetrieveSectorPayRatio(ctx, sector)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}
func (svc *service) CellPayRatio(ctx context.Context, cell string) (Chart, error) {
	const op errors.Op = "app/stats/service.CellPayRatio"

	chart, err := svc.repo.RetrieveCellPayRatio(ctx, cell)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}

func (svc *service) VillagePayRatio(ctx context.Context, village string) (Chart, error) {
	const op errors.Op = "app/stats/service.VillagePayRatio"

	chart, err := svc.repo.RetrieveVillagePayRatio(ctx, village)
	if err != nil {
		return Chart{}, errors.E(op, err)
	}

	return chart, nil
}
