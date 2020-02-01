package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Authenticate ...
func Authenticate(svc auth.Service) mux.MiddlewareFunc {
	const op errors.Op = "api/http/Authenticate"
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			token := strings.Split(r.Header.Get("Authorization"), "Bearer")[1]

			creds, err := svc.Identify(r.Context(), token)
			if err != nil {
				err = errors.E(op, err)
				w.WriteHeader(errors.Kind(err))
				w.Write([]byte(err.Error()))
			}

			ctx := auth.SetECredetialsInContext(r.Context(), &creds)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}
