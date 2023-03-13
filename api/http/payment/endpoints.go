package payment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
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

// ConfirmPull handles payment confirmation callback
func ConfirmPull(logger log.Entry, svc payment.Service) http.Handler {
	const op errors.Op = "api/http/payment/ConfirmPull"

	f := func(w http.ResponseWriter, r *http.Request) {

		callback := payment.Callback{}

		b, _ := ioutil.ReadAll(r.Body)
		logger.Debugf("request-body: '%s'", string(b))

		if err := json.Unmarshal(b, &callback); err != nil {
			err = errors.E(op, err, errors.KindBadRequest)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		err := svc.ConfirmPull(r.Context(), callback)
		if err != nil {
			err = errors.E(op, err)
			logger.SystemErr(err)
			encoding.EncodeError(w, errors.KindAlreadyExists, err)
			return
		}

		if err := encoding.Encode(w, http.StatusOK, []byte("payment confirmation received successfully")); err != nil {
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
