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
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handleAddProperty(svc, w, r)
		}).Methods("POST")

		router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleUpdateProperty(svc, w, r)
		}).Methods("PUT")

		router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleViewProperty(svc, w, r)
		}).Methods("GET")

		router.HandleFunc("/owners/{owner}", func(w http.ResponseWriter, r *http.Request) {
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
