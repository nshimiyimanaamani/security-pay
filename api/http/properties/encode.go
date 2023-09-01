package properties

import (
	"encoding/json"
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

func encode(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

func encodeErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var msg string

	if errors.Kind(err) == errors.KindUnexpected {
		msg = "internal server error"
	}

	msg = err.Error()

	w.WriteHeader(errors.Kind(err))
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
