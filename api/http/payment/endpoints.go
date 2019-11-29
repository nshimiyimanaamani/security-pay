package payment

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Initialize handles payment initialization
func Initialize(logger log.Entry, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		payment := payment.Payment{}

		if err := decode(r.Body, &payment); err != nil {
			EncodeError(w, err)
			logger.Errorf(err.Error())
			return
		}
		msg, err := svc.Initilize(payment)
		if err != nil {
			logger.Errorf(err.Error())
			EncodeError(w, err)
			return
		}
		if err := encode(w, msg); err != nil {
			logger.Errorf(err.Error())
			EncodeError(w, err)
		}

	}

	return http.HandlerFunc(f)
}

// Validate handles payment validatiob(confirmation)
func Validate(logger log.Entry, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		validation := payment.Validation{}
		res, err := svc.Validate(validation)
		if err != nil {
			logger.Errorf(err.Error())
		}
		if err := encode(w, res); err != nil {
			logger.Errorf(err.Error())
		}
	}

	return http.HandlerFunc(f)
}
