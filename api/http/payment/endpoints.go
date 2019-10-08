package payment

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/logger"
)

// Protocol adapts the payment service into an http.handler
type Protocol func(logger logger.Logger, svc payment.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Service payment.Service
	Logger  logger.Logger
}

// Initialize handles payment initialization
func Initialize(logger logger.Logger, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		payment := payment.Payment{}

		if err := decode(r.Body, &payment); err != nil {
			EncodeError(w, err)
			logger.Error(err.Error())
			return
		}
		msg, err := svc.Initilize(payment)
		if err != nil {
			logger.Error(err.Error())
			EncodeError(w, err)
			return
		}
		if err := encode(w, msg); err != nil {
			logger.Error(err.Error())
			EncodeError(w, err)
		}

	}

	return http.HandlerFunc(f)
}

// Validate handles payment validatiob(confirmation)
func Validate(logger logger.Logger, svc payment.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		validation := payment.Validation{}
		res, err := svc.Validate(validation)
		if err != nil {
			logger.Error(err.Error())
		}
		if err := encode(w, res); err != nil {
			logger.Error(err.Error())
		}
	}

	return http.HandlerFunc(f)
}
