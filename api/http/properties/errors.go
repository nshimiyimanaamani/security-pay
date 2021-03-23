package properties

import (
	"encoding/json"
	"strconv"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// parseTransErr parses http layer errors
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
