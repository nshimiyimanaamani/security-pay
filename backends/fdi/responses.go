package fdi

import validate "github.com/go-playground/validator/v10"

type statusResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data,omitempty"`
}

type authResponse struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func (res *authResponse) Validate() error {
	validator := validate.New()
	return validator.Struct(res)
}

type pullResponse struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
}

// Data ...
type Data struct {
	Token  string `json:"token,omitempty"`
	GwRef  string `json:"gwRef,omitempty"`
	State  string `json:"state,omitempty"`
	TrxRef string `json:"trxRef,omitempty"`
}
