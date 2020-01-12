package stats

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/stats"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// SectorPayRatio handles requests for the ratio payed/unpayed(pending) for a sector
func SectorPayRatio(logger log.Entry, svc stats.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)

		//var sector = vars["sector"]
	}

	return http.HandlerFunc(f)
}

// CellPayRatio handles requests for the ratio payed/unpayed(pending) for a sector
func CellPayRatio(logger log.Entry, svc stats.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)

		//var cell = vars["cell"]
	}

	return http.HandlerFunc(f)
}
