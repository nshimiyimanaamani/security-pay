package invoices

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Retrieve Handles invoice retrieval
func Retrieve(lgger log.Entry, svc invoices.Service) http.Handler {
	const op errors.Op = "api/http/invoices/Retrieve"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		property := vars["property"]

		months, err := strconv.ParseUint(vars["months"], 10, 64)
		if err != nil {
			err = errors.E(op, "could not parse months", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
		}

		res, err := svc.Retrieve(r.Context(), property, uint(months))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}
