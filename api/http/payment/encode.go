package payment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

var (
	contentType               = "application/json"
	errUnsupportedContentType = errors.New("unsupported content type")
)

func encode(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", contentType)
	return json.NewEncoder(w).Encode(v)
}

//EncodeError encodes the application error to the http api
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", contentType)

	var errMessage = newErrorMessage(err.Error())

	switch err {
	case transactions.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
	case transactions.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case transactions.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case transactions.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case users.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		errMessage = newErrorMessage(transactions.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		errMessage = newErrorMessage(transactions.ErrInvalidEntity.Error())
		w.WriteHeader(http.StatusBadRequest)

	default:
		switch err.(type) {
		case *json.SyntaxError:
			errMessage = newErrorMessage(transactions.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			errMessage = newErrorMessage(transactions.ErrInvalidEntity.Error())
			w.WriteHeader(http.StatusBadRequest)
		case *strconv.NumError:
			errMessage = newErrorMessage(transactions.ErrInvalidEntity.Error())
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
