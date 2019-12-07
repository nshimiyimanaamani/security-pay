package fdi

type statusResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data,omitempty"`
}

type authResponse struct {
	Status  string `json:"status"`
	Data    *Data  `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
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
