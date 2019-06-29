package http

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/logger"
)

//LoggingMiddleware logs every http request
func LoggingMiddleware(l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
