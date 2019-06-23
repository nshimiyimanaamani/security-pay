package http

import (
	"encoding/json"
	"github.com/rugwirobaker/paypack-backend/models"
	"net/http"
)

//UserRegisterEndpoint handles user registration
func (api *API) UserRegisterEndpoint(w http.ResponseWriter, r *http.Request) error {
	var err error

	var user models.User

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if err = user.Validate(); err != nil {
		return err
	}

	var id string

	id, err = api.Users.Register(user)
	if err != nil {
		return err
	}

	if err = encodeResponse(w, userRegisterResponse{ID: id}); err != nil {
		return err
	}
	return nil
}

//UserLoginEndpoint handles user login
func (api *API) UserLoginEndpoint(w http.ResponseWriter, r *http.Request) error {
	var err error

	var user models.User

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	if err = user.Validate(); err != nil {
		return err
	}

	var token string

	token, err = api.Users.Login(user)
	if err != nil {
		return err
	}

	if err = encodeResponse(w, userLoginResponse{Token: token}); err != nil {
		return err
	}
	return nil
}
