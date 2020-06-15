package properties

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(logger log.Entry, svc properties.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger        *log.Logger
	Service       properties.Service
	Authenticator auth.Service
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

	authenticator := middleware.Authenticate(opts.Logger, opts.Authenticator)

	r.Handle(RegisterPRoute, authenticator(LogEntryHandler(Register, opts))).
		Methods(http.MethodPost)

	r.Handle(RetrievePRoute, authenticator(LogEntryHandler(Retrieve, opts))).
		Methods(http.MethodGet)

	r.Handle(UpdatePRoute, authenticator(LogEntryHandler(Update, opts))).
		Methods(http.MethodPut)

	r.Handle(DeletePRoute, authenticator(LogEntryHandler(Delete, opts))).
		Methods(http.MethodDelete)

	r.Handle(ListPRoute, authenticator(LogEntryHandler(ListByCell, opts))).
		Methods(http.MethodGet).
		Queries("cell", "{cell}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, authenticator(LogEntryHandler(ListByOwner, opts))).
		Methods(http.MethodGet).
		Queries("owner", "{owner}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, authenticator(LogEntryHandler(ListBySector, opts))).
		Methods(http.MethodGet).
		Queries("sector", "{sector}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, authenticator(LogEntryHandler(ListByVillage, opts))).
		Methods(http.MethodGet).
		Queries("village", "{village}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, authenticator(LogEntryHandler(ListByRecorder, opts))).
		Methods(http.MethodGet).
		Queries("user", "{user}", "offset", "{offset}", "limit", "{limit}")

	//mobile temp
	r.Handle(MRetrievePRoute, MRetrieveProperty(opts.Logger, opts.Service)).Methods(http.MethodGet)
	r.Handle(MListPRoute, MListPropertyByCell(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("cell", "{cell}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(MListPRoute, MListPropertyByOwner(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("owner", "{owner}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(MListPRoute, MListPropertyBySector(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("sector", "{sector}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(MListPRoute, MListPropertyByVillage(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("village", "{village}", "offset", "{offset}", "limit", "{limit}")

}
