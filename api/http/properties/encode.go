package properties

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rugwirobaker/paypack-backend/app/users"

	transport "github.com/rugwirobaker/paypack-backend/api/http"
	"github.com/rugwirobaker/paypack-backend/app/properties"
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

	switch err {
	case properties.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
	case properties.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case properties.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case properties.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case users.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		errMessage = newErrorMessage(properties.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		errMessage = newErrorMessage(properties.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)

	default:
		switch err.(type) {
		case *json.SyntaxError:
			errMessage = newErrorMessage(properties.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			errMessage = newErrorMessage(properties.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *strconv.NumError:
			errMessage = newErrorMessage(properties.ErrInvalidEntity.Error())
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
