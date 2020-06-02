package users

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc users.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger        *log.Logger
	Service       users.Service
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

	// admins
	r.Handle(RegisterAdminRoute, authenticator(LogEntryHandler(RegisterAdmin, opts))).
		Methods(http.MethodPost)

	r.Handle(RetrieveAdminRoute, authenticator(LogEntryHandler(RetrieveAdmin, opts))).
		Methods(http.MethodGet)

	r.Handle(UpdateAdminCredsRoute, authenticator(LogEntryHandler(UpdateAdminCreds, opts))).
		Methods(http.MethodPut)

	r.Handle(RegisterAgentRoute, authenticator(LogEntryHandler(RegisterAgent, opts))).
		Methods(http.MethodPost)

	r.Handle(RetrieveAgentRoute, authenticator(LogEntryHandler(RetrieveAgent, opts))).
		Methods(http.MethodGet)

	//agents
	r.Handle(DeleteAgentRoute, authenticator(LogEntryHandler(DeleteAgent, opts))).
		Methods(http.MethodDelete)

	r.Handle(UpdateAgentRoute, authenticator(LogEntryHandler(UpdateAgentDetails, opts))).
		Methods(http.MethodPut)

	r.Handle(UpdateAgentCredsRoute, authenticator(LogEntryHandler(UpdateAgentsCreds, opts))).
		Methods(http.MethodPut)

	r.Handle(ListAgentsRoute, authenticator(LogEntryHandler(ListAgents, opts))).
		Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

	// developers
	r.Handle(RegisterDeveloperRoute, authenticator(LogEntryHandler(RegisterDeveloper, opts))).
		Methods(http.MethodPost)

	r.Handle(RetrieveDeveloperRoute, authenticator(LogEntryHandler(RetrieveDeveloper, opts))).
		Methods(http.MethodGet)

	r.Handle(DeleteDeveloperRoute, authenticator(LogEntryHandler(DeleteDeveloper, opts))).
		Methods(http.MethodDelete)

	r.Handle(UpdateDeveloperCredsRoute, authenticator(LogEntryHandler(UpdateDeveloperCreds, opts))).
		Methods(http.MethodPut)

	r.Handle(ListDevelopersRoute, authenticator(LogEntryHandler(ListDevelopers, opts))).
		Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

	//managers
	r.Handle(RegisterManagerRoute, authenticator(LogEntryHandler(RegisterManager, opts))).
		Methods(http.MethodPost)

	r.Handle(RetrieveManagerRoute, authenticator(LogEntryHandler(RetrieveManager, opts))).
		Methods(http.MethodGet)

	r.Handle(DeleteManagerRoute, authenticator(LogEntryHandler(DeleteManager, opts))).
		Methods(http.MethodDelete)

	r.Handle(UpdateManagerCredsRoute, authenticator(LogEntryHandler(UpdateManagerCreds, opts))).
		Methods(http.MethodPut)
	r.Handle(ListManagersRoute, authenticator(LogEntryHandler(ListManagers, opts))).
		Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")
}
