package properties

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// RegisterProperty handles property registration
func RegisterProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var property properties.Property

		if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
			EncodeError(w, err)
			return
		}

		if err := property.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		property, err := svc.RegisterProperty(ctx, property)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusCreated, property); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// UpdateProperty handles property update
func UpdateProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var property properties.Property
		if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
			EncodeError(w, err)
			return
		}

		if err := property.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		vars := mux.Vars(r)
		property.ID = vars["id"]

		err := svc.UpdateProperty(ctx, property)
		if err != nil {
			EncodeError(w, err)
			return
		}
		if err = EncodeResponse(w, http.StatusOK, property); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// RetrieveProperty handles property retrieval
func RetrieveProperty(lgger log.Entry, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		property, err := svc.RetrieveProperty(ctx, id)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, property); err != nil {
			EncodeError(w, err)
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
			EncodeError(w, err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		owner := vars["owner"]

		page, err := svc.ListPropertiesByOwner(ctx, owner, offset, limit)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, page); err != nil {
			EncodeError(w, err)
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
			EncodeError(w, err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		page, err := svc.ListPropertiesBySector(ctx, vars["sector"], offset, limit)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, page); err != nil {
			EncodeError(w, err)
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
			EncodeError(w, err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		page, err := svc.ListPropertiesByCell(ctx, vars["cell"], offset, limit)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, page); err != nil {
			EncodeError(w, err)
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
			EncodeError(w, err)
			return
		}
		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		page, err := svc.ListPropertiesByVillage(ctx, vars["village"], offset, limit)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, page); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}
