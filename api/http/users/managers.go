package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// RegisterManager ...
func RegisterManager(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RegisterManager"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Manager

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.RegisterManager(r.Context(), user)
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

// RetrieveManager ...
func RetrieveManager(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RetrieveManager"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["email"]

		res, err := svc.RetrieveManager(r.Context(), id)
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

// ListManagers ...
func ListManagers(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users. ListManagers"

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

		res, err := svc.ListManagers(r.Context(), offset, limit)
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

// UpdateManagerCreds ...
func UpdateManagerCreds(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateManagerCreds"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Manager

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		user.Email = vars["email"]

		err = svc.UpdateManagerCreds(r.Context(), user)
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

// DeleteManager ...
func DeleteManager(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAgentsDetails"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["email"]

		err := svc.DeleteManager(r.Context(), id)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res := map[string]string{"message": fmt.Sprintf("deleted user with %s", id)}

		if err := encode(w, http.StatusOK, res); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}
