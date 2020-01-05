package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// RegisterDeveloper ...
func RegisterDeveloper(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RegisterDeveloper"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Developer

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.RegisterDeveloper(r.Context(), user)
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

// RetrieveDeveloper ...
func RetrieveDeveloper(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RetrieveDeveloper"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["email"]

		res, err := svc.RetrieveDeveloper(r.Context(), id)
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

// ListDevelopers ...
func ListDevelopers(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.ListDevelopers"

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

		res, err := svc.ListDevelopers(r.Context(), offset, limit)
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

// UpdateDeveloperCreds ...
func UpdateDeveloperCreds(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateDeveloperCreds"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Developer

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		user.Email = vars["email"]

		err = svc.UpdateDeveloperCreds(r.Context(), user)
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

// DeleteDeveloper ...
func DeleteDeveloper(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAgentsDetails"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["email"]

		err := svc.DeleteDeveloper(r.Context(), id)
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
