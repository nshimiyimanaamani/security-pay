package ussd

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Callback ...
func Callback(lgger log.Entry, svc ussd.Service) http.Handler {
	const op errors.Op = "api/http/ussd.Callback"

	f := func(w http.ResponseWriter, r *http.Request) {

		var req ussd.SessionRequest

		err := encoding.Decode(r, &req)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encoding.EncodeErr(w, err)
			return
		}
		defer r.Body.Close()

		if err := req.Validate(); err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encoding.EncodeErr(w, err)
			return
		}
		if err := encoding.Encode(w, http.StatusOK, "OK"); err != nil {
			lgger.SystemErr(err)
			encoding.EncodeErr(w, err)
			return
		}
		lgger.Infof("%v", req)
	}
	return http.HandlerFunc(f)
}
