package properties

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/logger"
)

// Protocol adapts the feedback service into an http.handler
type Protocol func(logger logger.Logger, svc properties.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Service properties.Service
	Logger  logger.Logger
}

// RegisterProperty handles property registration
func RegisterProperty(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

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

		token := r.Header.Get("Authorization")

		property, err := svc.RegisterProperty(token, property)
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
func UpdateProperty(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

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
		token := r.Header.Get("Authorization")

		err := svc.UpdateProperty(token, property)
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
func RetrieveProperty(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		id := vars["id"]
		//token := r.Header.Get("Authorization")

		property, err := svc.RetrieveProperty(id)
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
func ListPropertyByOwner(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//token := r.Header.Get("Authorization")

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

		page, err := svc.ListPropertiesByOwner(owner, offset, limit)
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
func ListPropertyBySector(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//token := r.Header.Get("Authorization")

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

		page, err := svc.ListPropertiesBySector(vars["sector"], offset, limit)
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
func ListPropertyByCell(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//token := r.Header.Get("Authorization")

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

		page, err := svc.ListPropertiesByCell(vars["cell"], offset, limit)
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
func ListPropertyByVillage(logger logger.Logger, svc properties.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//token := r.Header.Get("Authorization")

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

		page, err := svc.ListPropertiesByVillage(vars["village"], offset, limit)
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
