package accounts

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Create handlers account creation
func Create(lgger log.Entry, svc accounts.Service) http.Handler {
	const op errors.Op = "api/http/accounts.Create"

	f := func(w http.ResponseWriter, r *http.Request) {
		var account accounts.Account

		err := Decode(r, &account)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		res, err := svc.Create(r.Context(), account)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		if err := encode(w, http.StatusCreated, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Retrieve handlers account creation
func Retrieve(lgger log.Entry, svc accounts.Service) http.Handler {
	const op errors.Op = "api/http/accounts.Retrieve"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		var id = vars["id"]

		res, err := svc.Retrieve(r.Context(), id)
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

// Update handlers account creation
func Update(lgger log.Entry, svc accounts.Service) http.Handler {
	const op errors.Op = "api/http/accounts.Update"

	f := func(w http.ResponseWriter, r *http.Request) {
		var account accounts.Account

		err := Decode(r, &account)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		account.ID = vars["id"]

		err = svc.Update(r.Context(), account)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("account[%s]: updated", account.ID)}); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// List handlers account creation
func List(lgger log.Entry, svc accounts.Service) http.Handler {
	const op errors.Op = "api/http/accounts.List"

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
