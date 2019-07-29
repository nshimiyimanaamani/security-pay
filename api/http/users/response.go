package users

import (
	"net/http"

	api "github.com/rugwirobaker/paypack-backend/api/http"
)

var (
	_ api.Response = (*userLoginResponse)(nil)
	_ api.Response = (*userRegisterResponse)(nil)
)

//UserRegisterResponse defines the http Response payload sent the Register endpoint
type userRegisterResponse struct {
	ID string `json:"id,omitempty"`
}

//UserLoginResponse defines the http Response payload sent to the login endpoint.
type userLoginResponse struct {
	Token string `json:"token,omitempty"`
}

func (res userRegisterResponse) Code() int {
	return http.StatusCreated
}

func (res userRegisterResponse) Headers() map[string]string {
	return map[string]string{}
}

func (res userRegisterResponse) Empty() bool {
	return res.ID == ""
}

func (res userLoginResponse) Code() int {
	return http.StatusCreated
}

func (res userLoginResponse) Headers() map[string]string {
	return map[string]string{}
}

func (res userLoginResponse) Empty() bool {
	return res.Token == ""
}
