package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Recover from server panics
func Recover() mux.MiddlewareFunc {
	const op errors.Op = "api/http/middleware/Recover"

	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					err := errors.E(op, err, errors.KindUnexpected)
					entry := log.EntryFromContext(r.Context())
					entry.SystemErr(err)
					w.WriteHeader(errors.Kind(err))
				}
			}()

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}
