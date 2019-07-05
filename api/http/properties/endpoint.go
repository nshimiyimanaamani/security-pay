package properties

import (
	"encoding/json"
	"net/http"

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

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		router.HandleFunc("/{owner}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		router.HandleFunc("/{sector}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		router.HandleFunc("/{cell}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		router.HandleFunc("/{village}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
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

func listPropertiesByOwner(svc properties.Service, w http.ResponseWriter, r *http.Request) {}

func listPropertiesBySector(svc properties.Service, w http.ResponseWriter, r *http.Request) {}

func listPropertiesByCell(svc properties.Service, w http.ResponseWriter, r *http.Request) {}

func listPropertiesByVillage(svc properties.Service, w http.ResponseWriter, r *http.Request) {}
