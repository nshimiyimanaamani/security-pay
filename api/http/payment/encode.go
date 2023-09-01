package payment

import (
	"encoding/json"
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

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
