package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Authenticate ...
func Authenticate(lgger log.Entry, svc auth.Service) mux.MiddlewareFunc {
	const op errors.Op = "api/http/Authenticate"
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				err := errors.E(op, "access denied: missing authorization token", errors.KindAccessDenied)
				lgger.SystemErr(err)
				encodeErr(w, errors.Kind(err), err)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")

			creds, err := svc.Identify(r.Context(), token)
			if err != nil {
				err = errors.E(op, err)
				lgger.SystemErr(err)
				encodeErr(w, errors.Kind(err), err)
				return
			}

			ctx := auth.SetECredetialsInContext(r.Context(), &creds)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}
