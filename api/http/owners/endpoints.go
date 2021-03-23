package owners

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Register handles owner registration
func Register(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var owner owners.Owner
		if err := json.NewDecoder(r.Body).Decode(&owner); err != nil {
			EncodeError(w, err)
			return
		}

		vars := mux.Vars(r)
		owner.ID = vars["id"]

		if err := owner.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		res, err := svc.Register(r.Context(), owner)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err := EncodeResponse(w, http.StatusCreated, owners.Owner{ID: res.ID}); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Retrieve handles owner retrieval
func Retrieve(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id := vars["id"]

		owner, err := svc.Retrieve(r.Context(), id)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, owner); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Update handles owner retrieval
func Update(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var err error

		if err = CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var owner owners.Owner
		if err = json.NewDecoder(r.Body).Decode(&owner); err != nil {
			EncodeError(w, err)
			return
		}

		vars := mux.Vars(r)
		owner.ID = vars["id"]

		if err = owner.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		if err = svc.Update(r.Context(), owner); err != nil {
			EncodeError(w, err)
			return
		}
		if err = EncodeResponse(w, http.StatusOK, owner); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// List handles multiple owners retrieval
func List(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		offset, err := strconv.ParseUint(vars["offset"], 10, 64)
		if err != nil {
			EncodeError(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 64)
		if err != nil {
			EncodeError(w, err)
			return
		}

		page, err := svc.List(r.Context(), offset, limit)
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

// Search handles owner search
func Search(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var owner owners.Owner

		vars := mux.Vars(r)

		owner.Fname = vars["fname"]
		owner.Lname = vars["lname"]
		owner.Phone = vars["phone"]

		res, err := svc.Search(r.Context(), owner)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, res); err != nil {
			EncodeError(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// RetrieveByPhone handles owner retrieval given phone number
func RetrieveByPhone(lgger log.Entry, svc owners.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		phone := vars["phone"]

		owner, err := svc.RetrieveByPhone(r.Context(), phone)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, http.StatusOK, owner); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}
