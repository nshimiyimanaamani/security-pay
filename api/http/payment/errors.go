package payment

import (
	"encoding/json"
	"strconv"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// errorMessage ..,
type errorMessage struct {
	Message string `json:"message"`
}

// newErrorMessage creates a new http error
func newErrorMessage(message string) errorMessage {
	return errorMessage{
		Message: message,
	}
}

func parseErr(op errors.Op, err error) error {
	switch err.(type) {
	case *json.SyntaxError:
		return errors.E(op, err, "invalid json syntax", errors.KindBadRequest)
	case *json.UnmarshalTypeError:
		return errors.E(op, err, "invalid json syntax", errors.KindBadRequest)
	case *strconv.NumError:
		return errors.E(op, err, "invalid request paramaters", errors.KindBadRequest)
	default:
		return err
	}
}
