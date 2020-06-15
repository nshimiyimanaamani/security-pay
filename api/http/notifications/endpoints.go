package notifications

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/core/notifications"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Send handles sms  notifications
func Send(lgger log.Entry, svc notifications.Service) http.Handler {
	const op errors.Op = "api/http/notifications/Send"
	f := func(w http.ResponseWriter, r *http.Request) {

		var req notifications.Payload

		err := decode(r, &req)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		if err := svc.Send(r.Context(), req); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusCreated, map[string]string{"message": "sms message was delivered"})
	}
	return http.HandlerFunc(f)
}
