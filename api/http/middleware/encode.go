package middleware

import (
	"encoding/json"
	"net/http"
)

// EncodeError ...
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	f := func(err error) map[string]string {
		return map[string]string{"error": err.Error()}
	}

	w.WriteHeader(http.StatusForbidden)

	json.NewEncoder(w).Encode(f(err))
}
