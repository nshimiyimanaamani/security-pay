package fdi

type statusResponse struct {
	Status string            `json:"status"`
	Data   map[string]string `json:"data"`
}

type authResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token   string `json:"token,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"data"`
}

//data[trxRef]
//data[gwRef]
//data[state]
type pullProcessResponse struct {
	Status string `json:"status"`
	Data   map[string]string
}

type pullResponse struct {
	Data struct {
		GwRef  string `json:"gwRef"`
		State  string `json:"state"`
		TrxRef string `json:"trxRef"`
	} `json:"data"`
	Status string `json:"status"`
}
