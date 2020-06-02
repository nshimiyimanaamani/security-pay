package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

func encode(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

//encodeError encodes the application error to the http api
func encodeErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
