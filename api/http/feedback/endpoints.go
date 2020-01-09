package feedback

import (
	"net/http"
	"strconv"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc feedback.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Service feedback.Service
	Logger  *log.Logger
}

// Recode handlers new feedback message submission
func Recode(lgger log.Entry, svc feedback.Service) http.Handler {
	const op errors.Op = "api/http/feedback/Record"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req feedback.Message

		err := decode(r, &req)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		res, err := svc.Record(ctx, &req)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusCreated, res)
	}

	return http.HandlerFunc(f)
}

// Update handles feedback updates
func Update(lgger log.Entry, svc feedback.Service) http.Handler {
	const op errors.Op = "api/http/feedback/Update"

	f := func(w http.ResponseWriter, r *http.Request) {

		var req feedback.Message

		err := decode(r, &req)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)

		req.ID = vars["id"]

		if err := svc.Update(r.Context(), req); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusOK, map[string]string{"message": "message updated"})
	}

	return http.HandlerFunc(f)
}

// Retrieve handles feedback entry retrieval
func Retrieve(lgger log.Entry, svc feedback.Service) http.Handler {
	const op errors.Op = "api/http/feedback/Retrieve"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		res, err := svc.Retrieve(ctx, id)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusOK, res)
	}

	return http.HandlerFunc(f)
}

// Delete handles feedback entry delete
func Delete(lgger log.Entry, svc feedback.Service) http.Handler {
	const op errors.Op = "api/http/feedback/Delete"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		err := svc.Delete(ctx, id)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusOK, map[string]string{"message": "message deleted"})
	}

	return http.HandlerFunc(f)
}

// List handles multiple message retrieval
func List(lgger log.Entry, svc feedback.Service) http.Handler {
	const op errors.Op = "api/http/feedback/List"

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
