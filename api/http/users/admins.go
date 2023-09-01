package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// RegisterAdmin ,,,
func RegisterAdmin(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RegisterAdmin"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Administrator

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.RegisterAdmin(r.Context(), user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusCreated, res); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// RetrieveAdmin ...
func RetrieveAdmin(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RetrieveAdmin"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["email"]

		res, err := svc.RetrieveAdmin(r.Context(), id)
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

// ListAdmins ...
func ListAdmins(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.ListAdmins"

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

		res, err := svc.ListAdmins(r.Context(), offset, limit)
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

// UpdateAdminCreds ...
func UpdateAdminCreds(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAdminCreds"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Administrator

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		user.Email = vars["email"]

		err = svc.UpdateAdminCreds(r.Context(), user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("user[%s]: updated", user.Email)}); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}
