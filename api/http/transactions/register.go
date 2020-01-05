package transactions

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(logger log.Entry, svc transactions.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger  *log.Logger
	Service transactions.Service
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

	r.Handle(RecordTransactionRoute, LogEntryHandler(Record, opts)).Methods(http.MethodPost)

	r.Handle(RetrieveTransactionRoute, LogEntryHandler(Retrieve, opts)).Methods(http.MethodGet)

	r.Handle(ListTransactionsRoute, LogEntryHandler(ListByProperty, opts)).Methods(http.MethodGet).
		Queries("property", "{property}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListTransactionsRoute, LogEntryHandler(ListByMethod, opts)).Methods(http.MethodGet).
		Queries("method", "{method}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListTransactionsRoute, LogEntryHandler(List, opts)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

	r.Handle(MListTransactionsRoute, LogEntryHandler(MListByProperty, opts)).Methods(http.MethodGet).
		Queries("property", "{property}")

}
