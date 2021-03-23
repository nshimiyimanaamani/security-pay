package auth

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Login handles user login
func Login(lgger log.Entry, svc auth.Service) http.Handler {
	const op errors.Op = "api/http/auth.Login"

	f := func(w http.ResponseWriter, r *http.Request) {
		var creds = auth.Credentials{}

		err := Decode(r, &creds)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}

		token, err := svc.Login(r.Context(), creds)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
		defer r.Body.Close()

		if err := encode(w, http.StatusOK, map[string]string{"token": token}); err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encodeErr(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Logout handles user login
func Logout(lgger log.Entry, svc auth.Service) http.Handler {
	const op errors.Op = "api/http/auth.Logout"
	f := func(w http.ResponseWriter, r *http.Request) {}

	return http.HandlerFunc(f)
}
