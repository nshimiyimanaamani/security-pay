package payment

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Initialize handles payment initialization
func Initialize(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Initialize"

	f := func(w http.ResponseWriter, r *http.Request) {
		tx := payment.Transaction{}

		err := encoding.Decode(r, &tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Initilize(r.Context(), tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
		if err := encoding.Encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Confirm handles payment confirmation callback
func Confirm(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Validate"

	f := func(w http.ResponseWriter, r *http.Request) {

		callback := payment.Callback{}

		if err := encoding.Decode(r, &callback); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		err := svc.Confirm(r.Context(), callback)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(f)
}
