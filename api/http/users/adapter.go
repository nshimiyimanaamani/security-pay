package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	transport "github.com/rugwirobaker/paypack-backend/api/http"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
)

//MakeAdapter takes a users service instance and returns a http handler
func MakeAdapter(mux *mux.Router) func(svc users.Service) {
	handler := func(svc users.Service) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var err error
			if err = transport.CheckContentType(r); err != nil {
				transport.EncodeError(w, err)
				return
			}

			var user models.User
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				transport.EncodeError(w, err)
				return
			}

			if err = user.Validate(); err != nil {
				transport.EncodeError(w, err)
				return
			}

			var id string

			id, err = svc.Register(user)
			if err != nil {
				transport.EncodeError(w, err)
				return
			}

			if err = transport.EncodeResponse(w, userRegisterResponse{ID: id}); err != nil {
				transport.EncodeError(w, err)
				return
			}
		}).Methods("POST")

		mux.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
			if err := transport.CheckContentType(r); err != nil {
				transport.EncodeError(w, err)
				return
			}
			var err error

			var user models.User

			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				transport.EncodeError(w, err)
				return
			}

			if err = user.Validate(); err != nil {
				transport.EncodeError(w, err)
				return
			}

			var token string

			token, err = svc.Login(user)
			if err != nil {
				transport.EncodeError(w, err)
				return
			}

			if err = transport.EncodeResponse(w, userLoginResponse{Token: token}); err != nil {
				transport.EncodeError(w, err)
			}

		}).Methods("POST")

		// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	}
	return handler
}
