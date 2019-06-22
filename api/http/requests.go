package http

import "github.com/rugwirobaker/paypack-backend/models"

//UserRegisterRequest defines the http Request payload sent the Register endpoint
type UserRegisterRequest struct {
	user models.User
}

func (req *UserRegisterRequest) validate() error {
	return req.user.Validate()
}

//UserLoginRequest defines the http Request payload sent to the login endpoint.
type UserLoginRequest struct {
	user models.User
}

func (req *UserLoginRequest) validate() error {
	return req.user.Validate()
}
