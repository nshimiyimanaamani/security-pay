package transactions

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHandlers ...
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(RecordTransactionRoute, Record(opts.Logger, opts.Service)).Methods(http.MethodPost)

	r.Handle(RetrieveTransactionRoute, Retrieve(opts.Logger, opts.Service)).Methods(http.MethodGet)

	r.Handle(ListTransactionsRoute, List(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")
}
