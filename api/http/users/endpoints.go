package users

import (
	"encoding/json"
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Register handles user registration
func Register(logger log.Entry, svc users.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var user users.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EncodeError(w, err)
			return
		}

		if err := user.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		var id string

		id, err := svc.Register(user)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err = EncodeResponse(w, userRegisterResponse{ID: id}); err != nil {
			EncodeError(w, err)
			return
		}
	}

	return http.HandlerFunc(f)
}

// Login handles user registration
func Login(logger log.Entry, svc users.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}
		var user users.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EncodeError(w, err)
			return
		}

		if err := user.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		var token string

		token, err := svc.Login(user)
		if err != nil {
			EncodeError(w, err)
			return
		}

		if err := EncodeResponse(w, userLoginResponse{Token: token}); err != nil {
			EncodeError(w, err)
		}
	}

	return http.HandlerFunc(f)
}
