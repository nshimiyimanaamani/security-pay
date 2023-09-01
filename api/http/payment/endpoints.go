package payment

import (
	"encoding/json"
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/api/http/encoding"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// Pull handles payment initialization
func Pull(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Pull"

	f := func(w http.ResponseWriter, r *http.Request) {

		tx := new(payment.TxRequest)

		err := encoding.Decode(r, &tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Pull(r.Context(), tx)
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

type msg struct {
	Message string `json:"message"`
}

// ProcessCallBack handles payment confirmation callback
func ProcessCallBack(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/ProcessCallBack"

	f := func(w http.ResponseWriter, r *http.Request) {

		callback := new(payment.Callback)

		dec := json.NewDecoder(r.Body)
		dec.UseNumber()
		if err := dec.Decode(callback); err != nil {
			err = errors.E(op, err, errors.KindBadRequest)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		err := svc.ProcessHook(r.Context(), *callback)
		if err != nil {
			out := msg{
				Message: err.Error(),
			}
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.Encode(w, errors.Kind(err), out)
			return
		}

		out := msg{
			Message: "payment confirmation received successfully",
		}

		if err := encoding.Encode(w, http.StatusOK, out); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Push handles payment initialization
func Push(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/Push"

	f := func(w http.ResponseWriter, r *http.Request) {

		tx := new(payment.TxRequest)

		err := encoding.Decode(r, &tx)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(errors.E(op, err))
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Push(r.Context(), tx)
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

// ConfirmPush handles payment confirmation callback
func ConfirmPush(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/ConfirmPush"

	f := func(w http.ResponseWriter, r *http.Request) {

		var callback payment.Callback

		if err := encoding.Decode(r, &callback); err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		err := svc.ConfirmPush(r.Context(), callback)
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
