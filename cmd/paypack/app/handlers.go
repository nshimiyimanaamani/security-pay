package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/accounts"
	"github.com/rugwirobaker/paypack-backend/api/http/auth"
	"github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/api/http/health"
	"github.com/rugwirobaker/paypack-backend/api/http/invoices"
	"github.com/rugwirobaker/paypack-backend/api/http/metrics"
	"github.com/rugwirobaker/paypack-backend/api/http/notifications"
	"github.com/rugwirobaker/paypack-backend/api/http/owners"
	"github.com/rugwirobaker/paypack-backend/api/http/payment"
	"github.com/rugwirobaker/paypack-backend/api/http/properties"
	"github.com/rugwirobaker/paypack-backend/api/http/transactions"
	"github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/api/http/version"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// HandlerOptions ...
type HandlerOptions struct {
	AccountsOptions *accounts.HandlerOpts
	AuthOptions     *auth.HandlerOpts
	FeedbackOptions *feedback.HandlerOpts
	NotifOptions    *notifications.HandlerOpts
	OwnersOptions   *owners.HandlerOpts
	PayOptions      *payment.HandlerOpts
	PropsOptions    *properties.HandlerOpts
	TransOptions    *transactions.HandlerOpts
	UsersOptions    *users.HandlerOpts
	InvoiceOptions  *invoices.HandlerOpts
	StatsOptions    *metrics.HandlerOpts
}

// NewHandlerOptions ...
func NewHandlerOptions(s *Services, lggr *log.Logger) *HandlerOptions {
	feedOpts := &feedback.HandlerOpts{
		Service: s.Feedback,
		Logger:  lggr,
	}
	proOpts := &properties.HandlerOpts{
		Service: s.Properties,
		Logger:  lggr,
	}

	ownersOpts := &owners.HandlerOpts{
		Service: s.Owners,
		Logger:  lggr,
	}
	paymentOpts := &payment.HandlerOpts{
		Service: s.Payment,
		Logger:  lggr,
	}

	transOpts := &transactions.HandlerOpts{
		Service: s.Transactions,
		Logger:  lggr,
	}

	usersOpts := &users.HandlerOpts{
		Service: s.Users,
		Logger:  lggr,
	}

	accountsOpts := &accounts.HandlerOpts{
		Service: s.Accounts,
		Logger:  lggr,
	}

	authOpts := &auth.HandlerOpts{
		Service: s.Auth,
		Logger:  lggr,
	}
	invOpts := &invoices.HandlerOpts{
		Service: s.Invoices,
		Logger:  lggr,
	}
	statsOpts := &metrics.HandlerOpts{
		Service: s.Stats,
		Logger:  lggr,
	}

	opts := &HandlerOptions{
		AuthOptions:     authOpts,
		AccountsOptions: accountsOpts,
		FeedbackOptions: feedOpts,
		OwnersOptions:   ownersOpts,
		PropsOptions:    proOpts,
		PayOptions:      paymentOpts,
		TransOptions:    transOpts,
		UsersOptions:    usersOpts,
		InvoiceOptions:  invOpts,
		StatsOptions:    statsOpts,
	}
	return opts
}

// Register registers all handlers
func Register(mux *mux.Router, opts *HandlerOptions) {
	if opts.FeedbackOptions == nil || opts.OwnersOptions == nil || opts.PropsOptions == nil {
		panic("absolutely unacceptable start server opts")
	}

	mux.HandleFunc("/healthz", health.Health).Methods(http.MethodGet)
	mux.HandleFunc("/version", version.Build).Methods(http.MethodGet)

	users.RegisterHandlers(mux, opts.UsersOptions)

	feedback.RegisterHandlers(mux, opts.FeedbackOptions)

	owners.RegisterHandlers(mux, opts.OwnersOptions)

	properties.RegisterHandlers(mux, opts.PropsOptions)

	payment.RegisterHandlers(mux, opts.PayOptions)

	transactions.RegisterHandlers(mux, opts.TransOptions)

	accounts.RegisterHandlers(mux, opts.AccountsOptions)

	auth.RegisterHandlers(mux, opts.AuthOptions)

	invoices.RegisterHandlers(mux, opts.InvoiceOptions)

	metrics.RegisterHandlers(mux, opts.StatsOptions)
}
