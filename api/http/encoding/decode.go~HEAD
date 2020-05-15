package encoding

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Decode requests
func Decode(r *http.Request, v interface{}) error {
	const op errors.Op = "api/http/properties.Decode"

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return errors.E(op, "invalid request: invalid content type", errors.KindUnsupportedContent)
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		switch err {
		case io.EOF:
			return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
		case io.ErrUnexpectedEOF:
			return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			// default:
			// 	switch err.(type) {
			// 	case *json.SyntaxError:
			// 		return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			// 	case *json.UnmarshalTypeError:
			// 		return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			// 	case *strconv.NumError:
			// 		return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			// 	default:
			// 		return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			// 	}
		}
	}
	return nil
}
