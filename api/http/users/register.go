package users

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc users.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger  *log.Logger
	Service users.Service
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

	r.Handle(RegisterAdminRoute, LogEntryHandler(RegisterAdmin, opts)).Methods(http.MethodPost)
	r.Handle(RetrieveAdminRoute, LogEntryHandler(RetrieveAdmin, opts)).Methods(http.MethodGet)
	r.Handle(UpdateAdminCredsRoute, LogEntryHandler(UpdateAdminCreds, opts)).Methods(http.MethodPut)

	r.Handle(RegisterAgentRoute, LogEntryHandler(RegisterAgent, opts)).Methods(http.MethodPost)
	r.Handle(RetrieveAgentRoute, LogEntryHandler(RetrieveAgent, opts)).Methods(http.MethodGet)
	r.Handle(UpdateAgentRoute, LogEntryHandler(UpdateAgentDetails, opts)).Methods(http.MethodPut)
	r.Handle(UpdateAgentCredsRoute, LogEntryHandler(UpdateAgentsCreds, opts)).Methods(http.MethodPut)
	r.Handle(ListAgentsRoute, LogEntryHandler(ListAgents, opts)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

	r.Handle(RegisterDeveloperRoute, LogEntryHandler(RegisterDeveloper, opts)).Methods(http.MethodPost)
	r.Handle(RetrieveDeveloperRoute, LogEntryHandler(RetrieveDeveloper, opts)).Methods(http.MethodGet)
	r.Handle(UpdateDeveloperCredsRoute, LogEntryHandler(UpdateDeveloperCreds, opts)).Methods(http.MethodPut)
	r.Handle(ListDevelopersRoute, LogEntryHandler(ListDevelopers, opts)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

	r.Handle(RegisterManagerRoute, LogEntryHandler(RegisterManager, opts)).Methods(http.MethodPost)
	r.Handle(RetrieveManagerRoute, LogEntryHandler(RetrieveManager, opts)).Methods(http.MethodGet)
	r.Handle(UpdateManagerCredsRoute, LogEntryHandler(UpdateManagerCreds, opts)).Methods(http.MethodPut)
	r.Handle(ListManagersRoute, LogEntryHandler(ListManagers, opts)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")
}
