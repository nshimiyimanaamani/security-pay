package payment

import (
	"encoding/json"
	"strings"

	"io"
	"net/http"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var (
	contentType = "application/json"
)

// VerifyContentType middleware checks content typ
func VerifyContentType(r *http.Request) error {
	const op errors.Op = ""
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return errors.E(op, "unsupported content type", errors.KindBadRequest)
	}
	return nil
}

func encodeRes(w io.Writer, i interface{}) error {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")

		return json.NewEncoder(w).Encode(i)
	}
	return nil
}

func encodeErr(w io.Writer, err error) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-Type", "application/json")

		data, _ := json.Marshal(map[string]string{"error": errors.KindText(err)})
		headered.Write(data)

	}
}
