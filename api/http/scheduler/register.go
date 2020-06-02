package scheduler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/scheduler"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// LogEntryHandler pulls a log entry from the request context. Thanks to the
// LogEntryMiddleware, we should have a log entry stored in the context for each
// request with request-specific fields. This will grab the entry and pass it to
// the protocol handlers
func LogEntryHandler(ph ProtocolHandler, opts *HandlerOpts) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ent := log.EntryFromContext(r.Context())
		handler := ph(ent, opts.Service)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

// ProtocolHandler adapts the scheduler service into an http.handler
type ProtocolHandler func(logger log.Entry, svc scheduler.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger        *log.Logger
	Service       scheduler.Service
	Authenticator auth.Service
}

// RegisterHandlers ...
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	authenticator := middleware.Authenticate(opts.Logger, opts.Authenticator)

	r.Handle(TasksRoute, authenticator(LogEntryHandler(Schedule, opts))).
		Methods(http.MethodGet)
}
