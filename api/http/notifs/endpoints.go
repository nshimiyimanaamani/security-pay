package notifs

import (
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
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
		// if _, err := svc.Send(r.Context(), sms); err != nil {
		// 	err = errors.E(op, err)
		// 	lgger.SystemErr(err)
		// 	encodeErr(w, errors.Kind(err), err)
		// 	return
		// }
		encode(w, http.StatusConflict, map[string]string{"message": "endpoint is under mentaince mode"})
	}
	return http.HandlerFunc(f)
}
