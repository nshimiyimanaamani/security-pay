package properties

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/logger"
)

// PageMetadata ...
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Property ...
type Property struct {
	ID         string  `json:"id"`
	Due        float64 `json:"due,string,omitempty"`
	OwnerID    string  `json:"owner_id"`
	OwnerFname string  `json:"owner_fname,omitempty"`
	OwnerLname string  `json:"owner_lname,omitempty"`
	OwnerPhone string  `json:"owner_phone,omitempty"`
	Sector     string  `json:"sector,omitempty"`
	Cell       string  `json:"cell,omitempty"`
	Village    string  `json:"village,omitempty"`
}

// PropertyPage represents a list of transaction.
type PropertyPage struct {
	PageMetadata
	Properties []Property
}

// MRetrieveProperty retrieves properties for mobile
func MRetrieveProperty(logger logger.Logger, svc properties.Service) http.Handler {
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

		response := Property{
			ID:         property.ID,
			Due:        property.Due,
			OwnerID:    property.Owner.ID,
			OwnerFname: property.Owner.Fname,
			OwnerLname: property.Owner.Lname,
			OwnerPhone: property.Owner.Phone,
			Sector:     property.Address.Sector,
			Cell:       property.Address.Cell,
			Village:    property.Address.Village,
		}

		if err = EncodeResponse(w, http.StatusOK, response); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MListPropertyByOwner handles property list by owner(for mobile)
func MListPropertyByOwner(logger logger.Logger, svc properties.Service) http.Handler {
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

		response := PropertyPage{
			PageMetadata: PageMetadata{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
		}

		for _, p := range page.Properties {
			property := Property{
				ID:         p.ID,
				Due:        p.Due,
				OwnerID:    p.Owner.ID,
				OwnerFname: p.Owner.Fname,
				OwnerLname: p.Owner.Lname,
				OwnerPhone: p.Owner.Phone,
				Sector:     p.Address.Sector,
				Cell:       p.Address.Cell,
				Village:    p.Address.Village,
			}
			response.Properties = append(response.Properties, property)
		}

		if err = EncodeResponse(w, http.StatusOK, response); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MListPropertyBySector handles property list by sector
func MListPropertyBySector(logger logger.Logger, svc properties.Service) http.Handler {
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

		response := PropertyPage{
			PageMetadata: PageMetadata{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
		}

		for _, p := range page.Properties {
			property := Property{
				ID:         p.ID,
				Due:        p.Due,
				OwnerID:    p.Owner.ID,
				OwnerFname: p.Owner.Fname,
				OwnerLname: p.Owner.Lname,
				OwnerPhone: p.Owner.Phone,
				Sector:     p.Address.Sector,
				Cell:       p.Address.Cell,
				Village:    p.Address.Village,
			}
			response.Properties = append(response.Properties, property)
		}

		if err = EncodeResponse(w, http.StatusOK, response); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MListPropertyByCell handles property list by cell(for mobile)
func MListPropertyByCell(logger logger.Logger, svc properties.Service) http.Handler {
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

		response := PropertyPage{
			PageMetadata: PageMetadata{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
		}

		for _, p := range page.Properties {
			property := Property{
				ID:         p.ID,
				Due:        p.Due,
				OwnerID:    p.Owner.ID,
				OwnerFname: p.Owner.Fname,
				OwnerLname: p.Owner.Lname,
				OwnerPhone: p.Owner.Phone,
				Sector:     p.Address.Sector,
				Cell:       p.Address.Cell,
				Village:    p.Address.Village,
			}
			response.Properties = append(response.Properties, property)

		}

		if err = EncodeResponse(w, http.StatusOK, response); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// MListPropertyByVillage handles property list by owner
func MListPropertyByVillage(logger logger.Logger, svc properties.Service) http.Handler {
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

		response := PropertyPage{
			PageMetadata: PageMetadata{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
		}

		for _, p := range page.Properties {
			property := Property{
				ID:         p.ID,
				Due:        p.Due,
				OwnerID:    p.Owner.ID,
				OwnerFname: p.Owner.Fname,
				OwnerLname: p.Owner.Lname,
				OwnerPhone: p.Owner.Phone,
				Sector:     p.Address.Sector,
				Cell:       p.Address.Cell,
				Village:    p.Address.Village,
			}
			response.Properties = append(response.Properties, property)
		}

		if err = EncodeResponse(w, http.StatusOK, response); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}
