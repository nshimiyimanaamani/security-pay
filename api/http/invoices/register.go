package invoices

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc invoices.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Service invoices.Service
	Logger  *log.Logger
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

// RegisterHandlers ....
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(RetrieveInvoicesRoute, LogEntryHandler(Retrieve, opts)).Methods(http.MethodGet).
		Queries("property", "{property}", "months", "{months}")

	r.Handle(MRetrieveAllInvoiceRoute, LogEntryHandler(MRetrieveAll, opts)).Methods(http.MethodGet).
		Queries("property", "{property}", "months", "{months}")

	r.Handle(MRetrievePendingInvoiceRoute, LogEntryHandler(MRetrievePending, opts)).Methods(http.MethodGet).
		Queries("property", "{property}", "months", "{months}")

	r.Handle(MRetrievePayedInvoiceRoute, LogEntryHandler(MRetrievePayed, opts)).Methods(http.MethodGet).
		Queries("property", "{property}", "months", "{months}")

}
