package metrics

import "context"

// Repository ...
type Repository interface {
	// FindSectorRatio returns the overall pending/payed ratio in a sector.
	FindSectorRatio(ctx context.Context, sector string) (Chart, error)

	// FindCellRatio returns the overall pending/payed ratio in a cell.
	FindCellRatio(ctx context.Context, cell string) (Chart, error)

	// FindVillageRatio returns the pending/payed ratio in a village
	FindVillageRatio(ctx context.Context, village string) (Chart, error)

	// ListSectorRatios returns the pending/payed ratio for all the cells in sector.
	ListSectorRatios(ctx context.Context, sector string) ([]Chart, error)

	// ListCellRatios returns the pending/payed ratio for all the villages in cell.
	ListCellRatios(ctx context.Context, sector string) ([]Chart, error)
}
