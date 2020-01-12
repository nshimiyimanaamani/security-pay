package stats

import "context"

// Service exposes stats usecases
type Service interface {
	SectorPayRatio(ctx context.Context, sector string) (Chart, error)
	CellPayRatio(ctx context.Context, cell string) (Chart, error)
	VillagePayRatio(ctx context.Context, village string) (Chart, error)
}
