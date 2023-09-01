package payment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/encoding"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/cast"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// PaymentReports returns the reports about paid and unpaid
func PaymentReports(logger log.Entry, svc payment.Repository) http.Handler {
	const op errors.Op = "api/http/payment/PaymentReports"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		status := cast.StringPointer((vars["status"]))
		sector := cast.StringPointer((vars["sector"]))
		cell := cast.StringPointer((vars["cell"]))
		village := cast.StringPointer((vars["village"]))
		from := cast.StringPointer((vars["from"]))
		to := cast.StringPointer((vars["to"]))

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
		creds := auth.CredentialsFromContext(r.Context())

		flt := &payment.Filters{
			Status:    status,
			Sector:    sector,
			Cell:      cell,
			Village:   village,
			From:      from,
			Namespace: &creds.Account,
			To:        to,
			Offset:    &offset,
			Limit:     &limit,
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

func TodayTransactions(logger log.Entry, svc payment.Repository) http.Handler {
	const op errors.Op = "api/http/payment/SectorPaymentMetrics"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		sector := cast.StringPointer((vars["sector"]))
		cell := cast.StringPointer(vars["cell"])
		village := cast.StringPointer(vars["village"])

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
		creds := auth.CredentialsFromContext(r.Context())
		flt := &payment.MetricFilters{

			Sector:    sector,
			Cell:      cell,
			Village:   village,
			Namespace: &creds.Account,
			Offset:    &offset,
			Limit:     &limit,
		}

		res, err := svc.TodayTransaction(r.Context(), flt)
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

func DailyTransactions(logger log.Entry, svc payment.Repository) http.Handler {

	const op errors.Op = "api/http/payment/DailyTransactions"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		sector := cast.StringPointer((vars["sector"]))
		cell := cast.StringPointer(vars["cell"])
		village := cast.StringPointer(vars["village"])
		from := cast.StringPointer(vars["from"])
		to := cast.StringPointer(vars["to"])

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

		if (cast.Uint64Pointer(offset) == nil || cast.Uint64Pointer(limit) == nil) || (offset == 0 && limit == 0) {
			offset = *cast.Uint64Pointer(0)
			limit = *cast.Uint64Pointer(20)
		}

		creds := auth.CredentialsFromContext(r.Context())
		flt := &payment.MetricFilters{

			Sector:    sector,
			Cell:      cell,
			Village:   village,
			From:      from,
			To:        to,
			Namespace: &creds.Account,
			Offset:    &offset,
			Limit:     &limit,
		}

		res, err := svc.ListDailyTransactions(r.Context(), flt)
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

// make summary transaction
func TodaySummary(logger log.Entry, svc payment.Repository) http.Handler {

	const op errors.Op = "api/http/payment/SummaryTransactions"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		sector := cast.StringPointer((vars["sector"]))
		cell := cast.StringPointer(vars["cell"])
		village := cast.StringPointer(vars["village"])
		date := cast.StringPointer((vars["date"]))

		creds := auth.CredentialsFromContext(r.Context())
		flt := &payment.MetricFilters{

			Sector:    sector,
			Cell:      cell,
			Village:   village,
			Date:      date,
			Namespace: &creds.Account,
		}

		res, err := svc.TodaySummary(r.Context(), flt)
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

// Unpaid houses
func UnpaidHouses(logger log.Entry, svc payment.Repository) http.Handler {

	const op errors.Op = "api/http/payment/UnpaidHouses"

	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		if vars["month"] == "" {
			now := time.Now()
			vars["month"] = strconv.Itoa(int(now.Month()))
		}

		month, err := strconv.Atoi(vars["month"])
		if err != nil {
			err = errors.E(op, err, "invalid month value", errors.KindBadRequest)
			logger.SystemErr(err)
			encodeErr(w, err)
			return
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

		if (cast.Uint64Pointer(offset) == nil || cast.Uint64Pointer(limit) == nil) || (offset == 0 && limit == 0) {
			offset = *cast.Uint64Pointer(0)
			limit = *cast.Uint64Pointer(20)
		}

		creds := auth.CredentialsFromContext(r.Context())
		flt := &payment.MetricFilters{

			Offset:    &offset,
			Limit:     &limit,
			Username:  &creds.Username,
			Namespace: &creds.Account,
			Month:     cast.Uint64Pointer(uint64(month)),
		}

		res, err := svc.UnpaidHouses(r.Context(), flt)
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
