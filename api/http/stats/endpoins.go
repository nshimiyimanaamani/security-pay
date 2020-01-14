package stats

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/stats"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// SectorPayRatio handles requests for the ratio payed/unpayed(pending) for a sector.
func SectorPayRatio(lgger log.Entry, svc stats.Service) http.Handler {
	const op errors.Op = "api/http/stats/SectorPayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var sector = vars["sector"]

		res, err := svc.SectorPayRatio(r.Context(), sector)
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
func CellPayRatio(lgger log.Entry, svc stats.Service) http.Handler {
	const op errors.Op = "api/http/stats/CellPayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["cell"]

		res, err := svc.CellPayRatio(r.Context(), cell)
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
func VillagePayRatio(lgger log.Entry, svc stats.Service) http.Handler {
	const op errors.Op = "api/http/stats/VillagePayRatio"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cell = vars["village"]

		res, err := svc.VillagePayRatio(r.Context(), cell)
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
