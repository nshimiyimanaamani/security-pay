package properties

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// RegisterProperty handles property registration
func RegisterProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		var property properties.Property

		err := Decode(r, &property)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.RegisterProperty(r.Context(), property)
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

// UpdateProperty handles property update
func UpdateProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		var property properties.Property

		err := Decode(r, &property)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		vars := mux.Vars(r)
		property.ID = vars["id"]

		res := svc.UpdateProperty(r.Context(), property)
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

// RetrieveProperty handles property retrieval
func RetrieveProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		id := vars["id"]

		res, err := svc.RetrieveProperty(r.Context(), id)
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

// ListPropertyByOwner handles property list by owner
func ListPropertyByOwner(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		owner := vars["owner"]

		res, err := svc.ListPropertiesByOwner(ctx, owner, offset, limit)
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

// ListPropertyBySector handles property list by sector
func ListPropertyBySector(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListPropertiesBySector(ctx, vars["sector"], offset, limit)
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

// ListPropertyByCell handles property list by cell
func ListPropertyByCell(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)

		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListPropertiesByCell(ctx, vars["cell"], offset, limit)
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

// ListPropertyByVillage handles property list by owner
func ListPropertyByVillage(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.ListPropertiesByVillage(ctx, vars["village"], offset, limit)
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
