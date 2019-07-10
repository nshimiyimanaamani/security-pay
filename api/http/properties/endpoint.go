package properties

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/properties"
)

// MakeEndpoint takes a proprerty service instance and returns a http handler
func MakeEndpoint(router *mux.Router) func(svc properties.Service) {
	handler := func(svc properties.Service) {
		//property handlers
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handleAddProperty(svc, w, r)
		}).Methods("POST")

		router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleUpdateProperty(svc, w, r)
		}).Methods("PUT")

		router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleViewProperty(svc, w, r)
		}).Methods("GET")

		router.HandleFunc("/owners/properties/{owner}", func(w http.ResponseWriter, r *http.Request) {
			handleListByOwner(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		router.HandleFunc("/sectors/{sector}", func(w http.ResponseWriter, r *http.Request) {
			handleListBySector(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		router.HandleFunc("/cells/{cell}", func(w http.ResponseWriter, r *http.Request) {
			handleListByCell(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		router.HandleFunc("/villages/{village}", func(w http.ResponseWriter, r *http.Request) {
			handleListByVillage(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		//owner handlers
		router.HandleFunc("/owners/", func(w http.ResponseWriter, r *http.Request) {
			handleCreateOwner(svc, w, r)
		}).Methods("POST")

		router.HandleFunc("/owners/", func(w http.ResponseWriter, r *http.Request) {
			handleListOwners(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		router.HandleFunc("/owners/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleViewOwner(svc, w, r)
		}).Methods("GET")

		router.HandleFunc("/owners/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleUpdateOwner(svc, w, r)
		}).Methods("PUT")

		router.HandleFunc("/owners/search/", func(w http.ResponseWriter, r *http.Request) {
			handleSearchOwner(svc, w, r)
		}).Queries("fname", "{fname}", "lname", "{lname}", "phone", "{phone}").Methods("GET")
	}
	return handler
}

func handleAddProperty(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = CheckContentType(r); err != nil {
		EncodeError(w, err)
		return
	}

	var property properties.Property

	if err = json.NewDecoder(r.Body).Decode(&property); err != nil {
		EncodeError(w, err)
		return
	}

	if err = property.Validate(); err != nil {
		EncodeError(w, err)
		return
	}

	property, err = svc.AddProperty(property)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := addPropertyRes{
		ID: property.ID,
	}
	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleUpdateProperty(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = CheckContentType(r); err != nil {
		EncodeError(w, err)
		return
	}

	var property properties.Property
	if err = json.NewDecoder(r.Body).Decode(&property); err != nil {
		EncodeError(w, err)
		return
	}

	vars := mux.Vars(r)
	property.ID = vars["id"]

	if err = property.Validate(); err != nil {
		EncodeError(w, err)
		return
	}

	err = svc.UpdateProperty(property)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := updatePropertyRes{
		ID:      property.ID,
		Owner:   property.Owner,
		Sector:  property.Sector,
		Cell:    property.Cell,
		Village: property.Village,
	}
	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleViewProperty(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	id := vars["id"]

	property, err := svc.ViewProperty(id)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := viewPropRes{
		ID:      property.ID,
		Owner:   property.Owner,
		Sector:  property.Sector,
		Cell:    property.Cell,
		Village: property.Village,
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleListByOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	page, err := svc.ListPropertiesByOwner(vars["owner"], offset, limit)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := propPageRes{
		pageRes: pageRes{
			Total:  page.Total,
			Offset: page.Offset,
			Limit:  page.Limit,
		},
		Properties: []viewPropRes{},
	}

	for _, property := range page.Properties {
		view := viewPropRes{
			ID:      property.ID,
			Owner:   property.Owner,
			Sector:  property.Sector,
			Cell:    property.Cell,
			Village: property.Village,
		}
		response.Properties = append(response.Properties, view)
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleListBySector(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	page, err := svc.ListPropertiesBySector(vars["sector"], offset, limit)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := propPageRes{
		pageRes: pageRes{
			Total:  page.Total,
			Offset: page.Offset,
			Limit:  page.Limit,
		},
		Properties: []viewPropRes{},
	}

	for _, property := range page.Properties {
		view := viewPropRes{
			ID:      property.ID,
			Owner:   property.Owner,
			Sector:  property.Sector,
			Cell:    property.Cell,
			Village: property.Village,
		}
		response.Properties = append(response.Properties, view)
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleListByCell(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	page, err := svc.ListPropertiesByCell(vars["cell"], offset, limit)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := propPageRes{
		pageRes: pageRes{
			Total:  page.Total,
			Offset: page.Offset,
			Limit:  page.Limit,
		},
		Properties: []viewPropRes{},
	}

	for _, property := range page.Properties {
		view := viewPropRes{
			ID:      property.ID,
			Owner:   property.Owner,
			Sector:  property.Sector,
			Cell:    property.Cell,
			Village: property.Village,
		}
		response.Properties = append(response.Properties, view)
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleListByVillage(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	page, err := svc.ListPropertiesByVillage(vars["village"], offset, limit)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := propPageRes{
		pageRes: pageRes{
			Total:  page.Total,
			Offset: page.Offset,
			Limit:  page.Limit,
		},
		Properties: []viewPropRes{},
	}

	for _, property := range page.Properties {
		view := viewPropRes{
			ID:      property.ID,
			Owner:   property.Owner,
			Sector:  property.Sector,
			Cell:    property.Cell,
			Village: property.Village,
		}
		response.Properties = append(response.Properties, view)
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleCreateOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = CheckContentType(r); err != nil {
		EncodeError(w, err)
		return
	}

	var owner properties.Owner

	if err = json.NewDecoder(r.Body).Decode(&owner); err != nil {
		EncodeError(w, err)
		return
	}

	if err = owner.Validate(); err != nil {
		EncodeError(w, err)
		return
	}

	owner.ID, err = svc.CreateOwner(owner)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := addPropertyRes{
		ID: owner.ID,
	}
	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleUpdateOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = CheckContentType(r); err != nil {
		EncodeError(w, err)
		return
	}

	var owner properties.Owner
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

	err = svc.UpdateOwner(owner)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := updateOwnerRes{
		ID:    owner.ID,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleViewOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	id := vars["id"]

	owner, err := svc.ViewOwner(id)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := viewOwnerRes{
		ID:    owner.ID,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleListOwners(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	page, err := svc.ListOwners(offset, limit)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := ownerPageRes{
		pageRes: pageRes{
			Total:  page.Total,
			Offset: page.Offset,
			Limit:  page.Limit,
		},
		Owners: []viewOwnerRes{},
	}

	for _, owner := range page.Owners {
		view := viewOwnerRes{
			ID:    owner.ID,
			Fname: owner.Fname,
			Lname: owner.Lname,
			Phone: owner.Phone,
		}
		response.Owners = append(response.Owners, view)
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}

func handleSearchOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	fname := vars["fname"]
	lname := vars["lname"]
	phone := vars["phone"]

	owner, err := svc.FindOwner(fname, lname, phone)
	if err != nil {
		EncodeError(w, err)
		return
	}

	response := viewOwnerRes{
		ID:    owner.ID,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	if err = EncodeResponse(w, response); err != nil {
		EncodeError(w, err)
		return
	}
}
