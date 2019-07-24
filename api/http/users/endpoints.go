package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

//MakeAdapter takes a users service instance and returns a http handler
func MakeAdapter(mux *mux.Router) func(svc users.Service) {
	handler := func(svc users.Service) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var err error
			if err = CheckContentType(r); err != nil {
				EncodeError(w, err)
				return
			}

			var user users.User
			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				EncodeError(w, err)
				return
			}

			if err = user.Validate(); err != nil {
				EncodeError(w, err)
				return
			}

			var id string

			id, err = svc.Register(user)
			if err != nil {
				EncodeError(w, err)
				return
			}

			if err = EncodeResponse(w, userRegisterResponse{ID: id}); err != nil {
				EncodeError(w, err)
				return
			}
		}).Methods("POST")

		mux.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
			if err := CheckContentType(r); err != nil {
				EncodeError(w, err)
				return
			}
			var err error

			var user users.User

			if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
				EncodeError(w, err)
				return
			}

			if err = user.Validate(); err != nil {
				EncodeError(w, err)
				return
			}

			var token string

			token, err = svc.Login(user)
			if err != nil {
				EncodeError(w, err)
				return
			}

			if err = EncodeResponse(w, userLoginResponse{Token: token}); err != nil {
				EncodeError(w, err)
			}

		}).Methods("POST")

		// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	}
	return handler
}
