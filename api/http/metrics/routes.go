package metrics

// metrics routes
const (
	SectorRatioRoute  = "/metrics/ratios/sectors/{sector}"
	CellRatioRoute    = "/metrics/ratios/cells/{cell}"
	VillageRatioRoute = "/metrics/ratios/villages/{village}"

	ListAllSectorRatiosRoute = "/metrics/ratios/sectors/all/{sector}"
	ListAllCellRatiosRoute   = "/metrics/ratios/cells/all/{cell}"

	SectorAccountBalance  = "/metrics/balance/sectors/{sector}"
	CellAccountBalance    = "/metrics/balance/cells/{cell}"
	VillageAccountBalance = "/metrics/balance/villages/{village}"
)
