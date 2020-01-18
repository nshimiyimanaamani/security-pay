package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/metrics"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// SectorPayRatio handles requests for the ratio payed/unpayed(pending) for a sector.
func SectorPayRatio(lgger log.Entry, svc metrics.Service) http.Handler {
	const op errors.Op = "api/http/metrics/SectorPayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		res, err := svc.FindSectorRatio(r.Context(), sector)
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

		res, err := svc.FindCellRatio(r.Context(), cell)
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

		res, err := svc.FindVillageRatio(r.Context(), cell)
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

		res, err := svc.ListAllSectorRatios(r.Context(), sector)
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

		res, err := svc.ListAllCellRatios(r.Context(), cell)
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
