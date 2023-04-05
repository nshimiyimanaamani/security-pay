package payment

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// ProtocolHandler adapts the payment service into an http.handler
type ProtocolHandler func(logger log.Entry, svc payment.Service) http.Handler
type RepoProtocolHandler func(logger log.Entry, repo payment.Repository) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger        *log.Logger
	Service       payment.Service
	Authenticator auth.Service
	Repository    payment.Repository
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

func RepoLogEntryHandler(ph RepoProtocolHandler, opts *HandlerOpts) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ent := log.EntryFromContext(r.Context())
		handler := ph(ent, opts.Repository)
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

	//temporary
	validator := middleware.ValidateRequestHeaders()

	r.Handle(ProcessDebitRoute, validator(LogEntryHandler(ProcessCallBack, opts))).Methods(http.MethodPost)
	r.Handle(DebitRoute, LogEntryHandler(Pull, opts)).Methods(http.MethodPost)

	r.Handle(ProcessCreditRoute, LogEntryHandler(ConfirmPush, opts)).Methods(http.MethodPost)
	r.Handle(CreditRoute, authenticator(LogEntryHandler(Push, opts))).Methods(http.MethodPost)
	r.Handle(PaymentReportsRoute, authenticator(RepoLogEntryHandler(PaymentReports, opts))).
		Methods(http.MethodGet).
		Queries("status", "{status}", "sector", "{sector}", "cell", "{cell}", "village", "{village}", "limit", "{limit}", "offset", "{offset}", "from", "{from}", "to", "{to}")

	r.Handle(TodayTransactionRoutes, authenticator(RepoLogEntryHandler(TodayTransactions, opts))).
		Methods(http.MethodGet).
		Queries("sector", "{sector}", "cell", "{cell}", "village", "{village}", "limit", "{limit}", "offset", "{offset}")

	r.Handle(DailyTransactionsRoutes, authenticator(RepoLogEntryHandler(DailyTransactions, opts))).Methods(http.MethodGet).
		Queries("sector", "{sector}", "cell", "{cell}", "village", "{village}", "from", "{from}", "to", "{to}", "limit", "{limit}", "offset", "{offset}")

	r.Handle(TodaySummaryRoute, authenticator(RepoLogEntryHandler(TodaySummary, opts))).Methods(http.MethodGet).
		Queries("sector", "{sector}", "cell", "{cell}", "village", "{village}")
}
