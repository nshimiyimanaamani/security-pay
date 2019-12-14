package accounts

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc accounts.Service) http.Handler

// HandlerOpts are the generic options
// for aProtocolHandler
type HandlerOpts struct {
	Logger  *log.Logger
	Service accounts.Service
}

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

// RegisterHandlers ...
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(CreateAccountRoute, LogEntryHandler(Create, opts)).Methods(http.MethodPost)
	r.Handle(UpdateAccountRoute, LogEntryHandler(Update, opts)).Methods(http.MethodPut)
	r.Handle(RetrieveAccountRoute, LogEntryHandler(Retrieve, opts)).Methods(http.MethodGet)

	r.Handle(ListAccountsRoute, LogEntryHandler(List, opts)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")
}
