package payment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/cast"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// PaymentReports returns the reports about paid and unpaid
func PaymentReports(logger log.Entry, svc payment.Repository) http.Handler {
	const op errors.Op = "api/http/payment/PaymentReports"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		status := cast.CastToString((vars["status"]))
		sector := cast.CastToString((vars["sector"]))
		cell := cast.CastToString((vars["cell"]))
		village := cast.CastToString((vars["village"]))
		from := cast.CastToString((vars["from"]))
		to := cast.CastToString((vars["to"]))
		var month int64
		if vars["month"] != "" {
			var err error
			month, err = strconv.ParseInt(vars["month"], 10, 32)
			if err != nil {
				err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
				logger.SystemErr(err)
				return
			}
		}

		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			err = errors.E(op, err, "invalid offset value", errors.KindBadRequest)
			logger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			err = errors.E(op, err, "invalid limit value", errors.KindBadRequest)
			logger.SystemErr(err)
			encodeErr(w, err)
			return
		}

		// default pagination settings if none are set
		if (cast.Uint64Pointer(offset) == nil || cast.Uint64Pointer(limit) == nil) || (offset == 0 && limit == 0) {
			offset = *cast.Uint64Pointer(0)
			limit = *cast.Uint64Pointer(20)
		}
		if month == 0 {
			// get the current month in number
			month = int64(time.Now().Month())
		}

		flt := &payment.Filters{
			Status:  status,
			Sector:  sector,
			Cell:    cell,
			Village: village,
			From:    from,
			To:      to,
			Month:   &month,
			Offset:  &offset,
			Limit:   &limit,
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
