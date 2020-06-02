package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/metrics"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(logger log.Entry, svc metrics.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger        *log.Logger
	Service       metrics.Service
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
	if opts == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	authenticator := middleware.Authenticate(opts.Logger, opts.Authenticator)

	//ratios
	r.Handle(SectorRatioRoute, authenticator(LogEntryHandler(SectorPayRatio, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

	r.Handle(CellRatioRoute, authenticator(LogEntryHandler(CellPayRatio, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

	r.Handle(VillageRatioRoute, authenticator(LogEntryHandler(VillagePayRatio, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

	r.Handle(ListAllCellRatiosRoute, authenticator(LogEntryHandler(ListAllCellRatios, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")
	r.Handle(ListAllSectorRatiosRoute, authenticator(LogEntryHandler(ListAllSectorRatios, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

	// balance
	r.Handle(SectorBalanceRoute, authenticator(LogEntryHandler(SectorBalance, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

	r.Handle(CellBalanceRoute, authenticator(LogEntryHandler(CellBalance, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")
	r.Handle(VillageBalanceRoute, authenticator(LogEntryHandler(VillageBalance, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")
	r.Handle(ListAllSectorBalancesRoute, authenticator(LogEntryHandler(ListAllSectorBalances, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")
	r.Handle(ListAllCellBalancesRoute, authenticator(LogEntryHandler(ListAllCellBalances, opts))).
		Methods(http.MethodGet).
		Queries("year", "{year}", "month", "{month}")

}
