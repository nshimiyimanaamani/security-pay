package metrics

// metrics routes
const (
	SectorRatioRoute         = "/metrics/ratios/sectors/{sector}"
	CellRatioRoute           = "/metrics/ratios/cells/{cell}"
	VillageRatioRoute        = "/metrics/ratios/villages/{village}"
	ListAllSectorRatiosRoute = "/metrics/ratios/sectors/all/{sector}"
	ListAllCellRatiosRoute   = "/metrics/ratios/cells/all/{cell}"

	SectorBalanceRoute         = "/metrics/balance/sectors/{sector}"
	CellBalanceRoute           = "/metrics/balance/cells/{cell}"
	VillageBalanceRoute        = "/metrics/balance/villages/{village}"
	ListAllSectorBalancesRoute = "/metrics/balance/sectors/all/{sector}"
	ListAllCellBalancesRoute   = "/metrics/balance/cells/all/{cell}"
)
