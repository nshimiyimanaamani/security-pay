package payment

import (
	"bytes"
	"io/ioutil"
	"net/http"
	lg "log"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Initialize handles payment initialization
func Initialize(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Initialize"

	f := func(w http.ResponseWriter, r *http.Request) {
		tx := payment.Transaction{}

		err := decode(r.Body, &tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encodeErr(w, errors.Kind(err), err)
			return
		}

		status, err := svc.Initilize(r.Context(), tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		if err := encodeRes(w, status); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Validate handles payment validatiob(confirmation)
func Validate(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Validate"

	f := func(w http.ResponseWriter, r *http.Request) {

		callback := payment.Callback{}

		
		buf, _ := ioutil.ReadAll(r.Body)

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		lg.Printf("body: %q", rdr1)

		if err := decode(r.Body, &callback); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		defer r.Body.Close()

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
