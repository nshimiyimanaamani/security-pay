package encoding

import (
	"encoding/json"
	"net/http"
)

func Encode(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

// EncodeError encodes the application error to the http api
func EncodeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
