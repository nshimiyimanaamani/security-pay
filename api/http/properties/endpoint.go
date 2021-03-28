package properties

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Register handles property registration
func Register(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/RegisterProperty"

	f := func(w http.ResponseWriter, r *http.Request) {

		var property properties.Property

		err := Decode(r, &property)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		res, err := svc.Register(r.Context(), property)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		defer r.Body.Close()

		if err := encode(w, http.StatusCreated, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Update handles property update
func Update(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/UpdateProperty"

	f := func(w http.ResponseWriter, r *http.Request) {

		var property properties.Property

		err := Decode(r, &property)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		vars := mux.Vars(r)
		property.ID = vars["id"]

		if err := svc.Update(r.Context(), property); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		res := map[string]string{"message": fmt.Sprintf("property: '%s' successfully updated", property.ID)}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Retrieve handles property retrieval
func Retrieve(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/RetrieveProperty"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		id := vars["id"]

		res, err := svc.Retrieve(r.Context(), id)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Delete handles property retrieval
func Delete(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/RetrieveProperty"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		id := vars["id"]

		err := svc.Delete(r.Context(), id)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
		encode(w, http.StatusOK, map[string]string{"message": "property deleted"})
	}
	return http.HandlerFunc(f)
}

// ListByOwner handles property list by owner
func ListByOwner(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertiesByOwner"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		owner := vars["owner"]

		res, err := svc.ListByOwner(ctx, owner, offset, limit)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListBySector handles property list by sector
func ListBySector(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertiesBySector"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		// creds := auth.CredentialsFromContext(ctx)

		// lgger.Warnf("username:%s | account:%s | role:%s",
		// 	creds.Username, creds.Account, creds.Role,
		// )

		res, err := svc.ListBySector(ctx, vars["sector"], offset, limit)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListByCell handles property list by cell
func ListByCell(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertiesByCell"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)

		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		res, err := svc.ListByCell(ctx, vars["cell"], offset, limit)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListByVillage handles property list by owner
func ListByVillage(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertiesByVillage"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		res, err := svc.ListByVillage(ctx, vars["village"], offset, limit)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// ListByRecorder handles property list by recorder
func ListByRecorder(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertiesByRecorder"

	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		res, err := svc.ListByRecorder(ctx, vars["user"], offset, limit)
		if err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		if err := encode(w, http.StatusOK, res); err != nil {
			err = parseErr(op, err)
			lgger.SystemErr(err)
			encodeErr(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}