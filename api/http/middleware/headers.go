package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ValidateRequestHeaders ...
func ValidateRequestHeaders() mux.MiddlewareFunc {
	const op errors.Op = "api/http/middleware/ValidateRequestHeaders"

	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
				err := errors.E(op, "invalid request: wrong content header", errors.KindUnsupportedContent)
				lgger := log.EntryFromContext(r.Context())
				lgger.SystemErr(err)
				w.WriteHeader(errors.Kind(err))
				encodeErr(w, errors.Kind(err), err)
			}
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}
