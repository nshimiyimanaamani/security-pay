package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	api "github.com/rugwirobaker/paypack-backend/api/http"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

//EncodeResponse encoded the application response to the http transport
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
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

//EncodeError encodes the application error to the http api
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", contentType)

	var errMessage = newErrorMessage(err.Error())

	switch err {
	case users.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
	case users.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case users.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case users.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		errMessage = newErrorMessage(users.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		errMessage = newErrorMessage(users.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)

	default:
		w.WriteHeader(http.StatusInternalServerError)

	}
	json.NewEncoder(w).Encode(errMessage)
}

//CheckContentType middleware checks content typ
func CheckContentType(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return errUnsupportedContentType
	}
	return nil
}
