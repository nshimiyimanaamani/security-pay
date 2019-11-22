package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/auth"
)

// Authenticate ...
func Authenticate(svc auth.Service) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			s := strings.Split(r.Header.Get("Authorization"), "Bearer")

			if err := svc.Authenticate(r.Context(), s[1]); err != nil {
				
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(f)
	}
}
