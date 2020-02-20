package notifications

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func decode(r *http.Request, v interface{}) error {
	const op errors.Op = "api/http/accounts.Decode"

	defer r.Body.Close()

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return errors.E(op, "invalid request: invalid content type", errors.KindUnsupportedContent)
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		switch err {
		case io.EOF:
			return errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
		default:
			switch err.(type) {
			case *json.SyntaxError:
				errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			case *json.UnmarshalTypeError:
				errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			default:
				errors.E(op, "invalid request: wrong data format", errors.KindBadRequest)
			}
		}
	}
	return nil
}
