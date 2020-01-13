package stats

import "context"

// Repository ...
type Repository interface {
	RetrieveSectorPayRatio(ctx context.Context, sector string) (Chart, error)
	RetrieveCellPayRatio(ctx context.Context, cell string) (Chart, error)
	RetrieveVillagePayRatio(ctx context.Context, village string) (Chart, error)
}
