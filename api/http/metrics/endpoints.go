package metrics

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nshimiyimanaamani/paypack-backend/core/metrics"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// SectorPayRatio handles requests for the ratio payed/unpayed(pending) for a sector.
func SectorPayRatio(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/SectorPayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindSectorRatio(r.Context(), sector, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// CellPayRatio handles requests for the ratio payed/unpayed(pending) for a cell.
func CellPayRatio(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/CellPayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["cell"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindCellRatio(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// VillagePayRatio handles requests for the ratio payed/unpayed(pending) for a village.
func VillagePayRatio(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/VillagePayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["village"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindVillageRatio(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListAllSectorRatios handles requests for the ratio payed/unpayed(pending) for a sector.
func ListAllSectorRatios(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/ListAllSectorRatios"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListAllSectorRatios(r.Context(), sector, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListAllCellRatios handles requests for the ratio payed/unpayed(pending) for a sector.
func ListAllCellRatios(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/ListAllCellRatios"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["cell"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListAllCellRatios(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// SectorBalance handles requests for the balance for a sector.
func SectorBalance(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/SectorBalance"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindSectorBalance(r.Context(), sector, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// CellBalance handles requests for a cell's balance.
func CellBalance(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/CellBalance"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["cell"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindCellBalance(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// VillageBalance handles requests for the ratio payed/unpayed(pending) for a village.
func VillageBalance(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/VillageBalance "

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["village"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.FindVillageBalance(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListAllSectorBalances  handles requests for the ratio payed/unpayed(pending) for a sector.
func ListAllSectorBalances(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/ListAllSectorBalances"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListAllSectorBalances(r.Context(), sector, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListAllCellBalances handles requests for the ratio payed/unpayed(pending) for a sector.
func ListAllCellBalances(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/ListAllCellBalances"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["cell"]

		year, err := strconv.ParseUint(vars["year"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid year value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		month, err := strconv.ParseUint(vars["month"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListAllCellBalances(r.Context(), cell, uint(year), uint(month))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}
