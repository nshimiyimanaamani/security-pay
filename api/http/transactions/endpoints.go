package transactions

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// Record handles transaction record
func Record(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.Record"

	f := func(w http.ResponseWriter, r *http.Request) {

		var transaction transactions.Transaction

		err := Decode(r, &transaction)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := transaction.Validate(); err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Record(r.Context(), transaction)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusCreated, res); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Retrieve handles transaction retrieval
func Retrieve(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.Retrieve"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["id"]

		res, err := svc.Retrieve(r.Context(), id)
		if err != nil {
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

// List handles transaction list
func List(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.List"
	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		offset, err := strconv.ParseUint(vars["offset"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid offset value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid limit value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.List(r.Context(), offset, limit)
		if err != nil {
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

// ListByProperty handles transactions retrieval given house code
func ListByProperty(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.ListByProperty"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var property = vars["property"]

		offset, err := strconv.ParseUint(vars["offset"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid offset value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid limit value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListByProperty(r.Context(), property, offset, limit)
		if err != nil {
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

// ListByMethod handles transactions retrieval given the transaction method
func ListByMethod(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.ListByMethod"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var method = vars["method"]

		offset, err := strconv.ParseUint(vars["offset"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid offset value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 64)
		if err != nil {
			err = errors.E(op, err, "invalid limit value", errors.KindBadRequest)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListByMethod(r.Context(), method, offset, limit)
		if err != nil {
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

// MListByProperty ...
func MListByProperty(lgger log.Entry, svc transactions.Service) http.Handler {
	const op errors.Op = "api/http/transactions.MListByProperty"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var property = vars["property"]

		page, err := svc.ListByPropertyR(r.Context(), property)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusOK, page.Transactions); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}
