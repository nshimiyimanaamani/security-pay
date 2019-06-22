package http

import (
	"encoding/json"
	"net/http"

	"github.com/rugwirobaker/paypack-backend/api"
)

var (
	contentType = "application/json"
)

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(api.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", contentType)
}
