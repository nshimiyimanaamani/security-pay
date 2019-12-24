package properties

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// RegisterProperty handles property registration
func RegisterProperty(lgger log.Entry, svc properties.Service) http.Handler {
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

		res, err := svc.RegisterProperty(r.Context(), property)
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

// UpdateProperty handles property update
func UpdateProperty(lgger log.Entry, svc properties.Service) http.Handler {
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

		if err := svc.UpdateProperty(r.Context(), property); err != nil {
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

// RetrieveProperty handles property retrieval
func RetrieveProperty(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/RetrieveProperty"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		id := vars["id"]

		res, err := svc.RetrieveProperty(r.Context(), id)
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

// ListPropertyByOwner handles property list by owner
func ListPropertyByOwner(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertyByOwner"

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

		res, err := svc.ListPropertiesByOwner(ctx, owner, offset, limit)
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

// ListPropertyBySector handles property list by sector
func ListPropertyBySector(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertyBySector"

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

		res, err := svc.ListPropertiesBySector(ctx, vars["sector"], offset, limit)
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

// ListPropertyByCell handles property list by cell
func ListPropertyByCell(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertyByCell"

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

		res, err := svc.ListPropertiesByCell(ctx, vars["cell"], offset, limit)
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

// ListPropertyByVillage handles property list by owner
func ListPropertyByVillage(lgger log.Entry, svc properties.Service) http.Handler {
	const op errors.Op = "api/http/properties/ListPropertyByVillage"

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

		res, err := svc.ListPropertiesByVillage(ctx, vars["village"], offset, limit)
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
