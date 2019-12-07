package payment

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Initialize handles payment initialization
func Initialize(logger log.Entry, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		const op errors.Op = "api.http.Payment.Initialize"
		ctx := r.Context()
		tx := payment.Transaction{}

		err := decode(r.Body, &tx)
		if err != nil {
			logger.SystemErr(errors.E(op, err))
			w.WriteHeader(http.StatusInternalServerError)
			encodeErr(w, err)
			return
		}

		status, err := svc.Initilize(ctx, tx)
		if err != nil {
			severityLevel := errors.Expect(err, errors.KindNotFound)
			err = errors.E(op, err, severityLevel)
			logger.SystemErr(err)
			w.WriteHeader(errors.Kind(err))
			encodeErr(w, err)
			return
		}
		if err := encodeRes(w, status); err != nil {
			logger.SystemErr(errors.E(op, err))
			w.WriteHeader(http.StatusInternalServerError)
			encodeErr(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Validate handles payment validatiob(confirmation)
func Validate(logger log.Entry, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		callback := payment.Callback{}
		if err := decode(r.Body, &callback); err != nil {
			encodeErr(w, err)
			logger.SystemErr(err)
			return
		}
		err := svc.Confirm(r.Context(), callback)
		if err != nil {
			encodeErr(w, err)
			logger.SystemErr(err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(f)
}
