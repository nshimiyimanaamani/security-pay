package owners

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rugwirobaker/paypack-backend/app/users"

	"github.com/rugwirobaker/paypack-backend/app/owners"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

//EncodeResponse encoded the application response to the http ts
func EncodeResponse(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

//EncodeError encodes the application error to the http api
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", contentType)

	var errMessage = newErrorMessage(err.Error())

	switch err {
	case owners.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
	case owners.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case owners.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case owners.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case users.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		errMessage = newErrorMessage(owners.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		errMessage = newErrorMessage(owners.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)

	default:
		switch err.(type) {
		case *json.SyntaxError:
			errMessage = newErrorMessage(owners.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			errMessage = newErrorMessage(owners.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *strconv.NumError:
			errMessage = newErrorMessage(owners.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	json.NewEncoder(w).Encode(errMessage)
}

//CheckContentType middleware checks content typ
func CheckContentType(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		//logger.Warn("Invalid or missing content type.")
		return errUnsupportedContentType
	}
	return nil
}
