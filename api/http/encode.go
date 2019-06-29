package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	
	"github.com/rugwirobaker/paypack-backend/models"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

//EncodeResponse encoded the application response to the http transport
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(Response); ok {
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

	switch err {
	case models.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
	case models.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case models.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case models.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)

	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//CheckContentType middleware checks content typ
func CheckContentType(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		//logger.Warn("Invalid or missing content type.")
		return errUnsupportedContentType
	}
	return nil
}
