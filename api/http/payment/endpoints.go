package payment

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Initialize handles payment initialization
func Initialize(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Initialize"

	f := func(w http.ResponseWriter, r *http.Request) {
		tx := payment.Transaction{}

		err := Decode(r, &tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encodeErr(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Initilize(r.Context(), tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		if err := encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encodeErr(w, errors.Kind(err), err)
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

		if err := Decode(r, &callback); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		err := svc.Confirm(r.Context(), callback)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(f)
}
