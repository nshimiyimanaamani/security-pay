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

// RegisterAgent ...
func RegisterAgent(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RegisterAgent"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Agent

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.RegisterAgent(r.Context(), user)
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

// RetrieveAgent ...
func RetrieveAgent(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.RetrieveAgent"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		var id = vars["phone"]

		res, err := svc.RetrieveAgent(r.Context(), id)
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

// ListAgents ...
func ListAgents(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.ListAgents"

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

		res, err := svc.ListAgents(r.Context(), offset, limit)
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

// UpdateAgentsCreds ...
func UpdateAgentsCreds(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAgentsCreds"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Agent

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		user.Telephone = vars["phone"]

		err = svc.UpdateAgentCreds(r.Context(), user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("user[%s]: updated", user.Telephone)}); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// UpdateAgentDetails ...
func UpdateAgentDetails(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAgentsDetails"

	f := func(w http.ResponseWriter, r *http.Request) {
		var user users.Agent

		err := Decode(r, &user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		user.Telephone = vars["phone"]

		err = svc.UpdateAgent(r.Context(), user)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("user[%s]: updated", user.Telephone)}); err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// DeleteAgent ...
func DeleteAgent(lgger log.Entry, svc users.Service) http.Handler {
	const op errors.Op = "api/http/users.UpdateAgentsDetails"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var id = vars["phone"]

		err := svc.DeleteAgent(r.Context(), id)
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
