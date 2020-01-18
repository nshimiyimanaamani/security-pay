package metrics

// metrics routes
const (
	SectorRatioRoute  = "/metrics/sectors/{sector}"
	CellRatioRoute    = "/metrics/cells/{cell}"
	VillageRatioRoute = "/metrics/villages/{village}"

	ListAllSectorRatiosRoute = "/metrics/sectors/all/{sector}"
	ListAllCellRatiosRoute   = "/metrics/cells/all/{cell}"

	CellsAccountBalance    = "/metrics/balance/sectors/{sector}"
	VillagesAccountBalance = "/metrics/balance/cells/{cell}"
)
