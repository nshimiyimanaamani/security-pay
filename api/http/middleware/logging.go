package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/logger"
)

//Log logs every http request
func Log(l logger.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {}

		return http.HandlerFunc(f)
	}
}
