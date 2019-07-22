package version

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	transport "github.com/rugwirobaker/paypack-backend/api/http"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

//EncodeResponse encoded the application response to the http transport
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(transport.Response); ok {
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

//EncodeError encodes the application error to the http api
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", contentType)

	var errMessage = newErrorMessage(err.Error())

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errMessage)
}

//CheckContentType middleware checks content typ
func CheckContentType(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return errUnsupportedContentType
	}
	return nil
}
