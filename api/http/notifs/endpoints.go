package notifs

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Send handles sms  notifs
func Send(lgger log.Entry, svc notifs.Service) http.Handler {
	const op errors.Op = "api/http/notifs/Send"
	f := func(w http.ResponseWriter, r *http.Request) {

		var sms notifs.Notification

		creds := auth.CredentialsFromContext(r.Context())
		if creds == nil {
			err := errors.E(op, "not logged in", errors.KindAccessDenied)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		sms.Sender = creds.Account

		err := decode(r, &sms)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		if _, err := svc.Send(r.Context(), sms); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		encode(w, http.StatusCreated, map[string]string{"message": "sms message was delivered"})
	}
	return http.HandlerFunc(f)
}
