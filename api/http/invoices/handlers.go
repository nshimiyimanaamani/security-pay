package invoices

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
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

		res, err := svc.RetrieveAll(r.Context(), property, uint(months))
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

// MRetrieveAll handles invoice request without metadata
func MRetrieveAll(lgger log.Entry, svc invoices.Service) http.Handler {
	const op errors.Op = "api/http/invoices/OnMobileRetrieve"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		property := vars["property"]

		months, err := strconv.ParseUint(vars["months"], 10, 64)
		if err != nil {
			err = errors.E(op, "could not parse months", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
		}

		page, err := svc.RetrieveAll(r.Context(), property, uint(months))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, page.Invoices); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MRetrievePending handles invoice request without metadata
func MRetrievePending(lgger log.Entry, svc invoices.Service) http.Handler {
	const op errors.Op = "api/http/invoices/MRetrievePending"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		property := vars["property"]

		months, err := strconv.ParseUint(vars["months"], 10, 64)
		if err != nil {
			err = errors.E(op, "could not parse months", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
		}

		page, err := svc.RetrievePending(r.Context(), property, uint(months))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, page.Invoices); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MRetrievePayed handles invoice request without metadata
func MRetrievePayed(lgger log.Entry, svc invoices.Service) http.Handler {
	const op errors.Op = "api/http/invoices/MRetrievePayed"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		property := vars["property"]

		months, err := strconv.ParseUint(vars["months"], 10, 64)
		if err != nil {
			err = errors.E(op, "could not parse months", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
		}

		page, err := svc.RetrievePayed(r.Context(), property, uint(months))
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, page.Invoices); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}
