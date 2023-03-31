package payment

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// PaymentReports returns the reports about paid and unpaid
func PaymentReports(logger log.Entry, svc payment.Repository) http.Handler {
	const op errors.Op = "api/http/payment/PaymentReports"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		status := vars["status"]

		// !!! TODO: Add the missing filters and other things too below
		// sector := vars["sector"]
		// cell := vars["cell"]
		// village := vars["village"]

		// offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		// if err != nil {
		// 	err = parseErr(op, err)
		// 	logger.SystemErr(err)
		// 	encodeErr(w, err)
		// 	return
		// }

		// limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		// if err != nil {
		// 	err = parseErr(op, err)
		// 	logger.SystemErr(err)
		// 	encodeErr(w, err)
		// 	return
		// }

		flt := &payment.Filters{
			Status: &status,
		}

		res, err := svc.List(r.Context(), flt)
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
