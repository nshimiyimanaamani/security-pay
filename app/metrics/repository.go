package metrics

import "context"

// Repository ...
type Repository interface {
	CountMetrics
	BalanceMetrics
}

// BalanceMetrics gives access to payment amount aggregations.
type BalanceMetrics interface {
	// FindSectorRatio returns the overall pending/payed ratio in a sector.
	FindSectorBalance(ctx context.Context, sector string, y, m uint) (Chart, error)

	// FindCellRatio returns the overall pending/payed ratio in a cell.
	FindCellBalance(ctx context.Context, cell string, y, m uint) (Chart, error)

	// FindVillageRatio returns the pending/payed ratio in a village
	FindVillageBalance(ctx context.Context, village string, y, m uint) (Chart, error)

	// ListSectorRatios returns the pending/payed ratio for all the cells in sector.
	ListSectorBalances(ctx context.Context, sector string, y, m uint) ([]Chart, error)

	// ListCellRatios returns the pending/payed ratio for all the villages in cell.
	ListCellBalances(ctx context.Context, sector string, y, m uint) ([]Chart, error)
}

// CountMetrics gives access to payments count aggregations.
type CountMetrics interface {
	// FindSectorRatio returns the overall pending/payed ratio in a sector.
	FindSectorRatio(ctx context.Context, sector string, y, m uint) (Chart, error)

	// FindCellRatio returns the overall pending/payed ratio in a cell.
	FindCellRatio(ctx context.Context, cell string, y, m uint) (Chart, error)

	// FindVillageRatio returns the pending/payed ratio in a village
	FindVillageRatio(ctx context.Context, village string, y, m uint) (Chart, error)

	// ListSectorRatios returns the pending/payed ratio for all the cells in sector.
	ListSectorRatios(ctx context.Context, sector string, y, m uint) ([]Chart, error)

	// ListCellRatios returns the pending/payed ratio for all the villages in cell.
	ListCellRatios(ctx context.Context, sector string, y, m uint) ([]Chart, error)
}
