package app

import (
	"net/http"

	"github.com/gorilla/mux"
	accountsEndpoints "github.com/rugwirobaker/paypack-backend/api/http/accounts"
	authEndpoints "github.com/rugwirobaker/paypack-backend/api/http/auth"
	feedbackEndpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/api/http/health"
	invoiceEndpoints "github.com/rugwirobaker/paypack-backend/api/http/invoices"
	ownersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/owners"
	paymentEndpoints "github.com/rugwirobaker/paypack-backend/api/http/payment"
	propertiesEndpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	metricsEndpoints "github.com/rugwirobaker/paypack-backend/api/http/metrics"
	transactionsEndpoints "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	usersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/api/http/version"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// HandlerOptions ...
type HandlerOptions struct {
	AccountsOptions *accountsEndpoints.HandlerOpts
	AuthOptions     *authEndpoints.HandlerOpts
	FeedbackOptions *feedbackEndpoints.HandlerOpts
	OwnersOptions   *ownersEndpoints.HandlerOpts
	PayOptions      *paymentEndpoints.HandlerOpts
	PropsOptions    *propertiesEndpoints.HandlerOpts
	TransOptions    *transactionsEndpoints.HandlerOpts
	UsersOptions    *usersEndpoints.HandlerOpts
	InvoiceOptions  *invoiceEndpoints.HandlerOpts
	StatsOptions    *metricsEndpoints.HandlerOpts
}

// NewHandlerOptions ...
func NewHandlerOptions(s *Services, lggr *log.Logger) *HandlerOptions {
	feedOpts := &feedbackEndpoints.HandlerOpts{
		Service: s.Feedback,
		Logger:  lggr,
	}
	proOpts := &propertiesEndpoints.HandlerOpts{
		Service: s.Properties,
		Logger:  lggr,
	}

	ownersOpts := &ownersEndpoints.HandlerOpts{
		Service: s.Owners,
		Logger:  lggr,
	}
	paymentOpts := &paymentEndpoints.HandlerOpts{
		Service: s.Payment,
		Logger:  lggr,
	}

	transOpts := &transactionsEndpoints.HandlerOpts{
		Service: s.Transactions,
		Logger:  lggr,
	}

	usersOpts := &usersEndpoints.HandlerOpts{
		Service: s.Users,
		Logger:  lggr,
	}

	accountsOpts := &accountsEndpoints.HandlerOpts{
		Service: s.Accounts,
		Logger:  lggr,
	}

	authOpts := &authEndpoints.HandlerOpts{
		Service: s.Auth,
		Logger:  lggr,
	}
	invOpts := &invoiceEndpoints.HandlerOpts{
		Service: s.Invoices,
		Logger:  lggr,
	}
	statsOpts := &metricsEndpoints.HandlerOpts{
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

	usersEndpoints.RegisterHandlers(mux, opts.UsersOptions)

	feedbackEndpoints.RegisterHandlers(mux, opts.FeedbackOptions)

	ownersEndpoints.RegisterHandlers(mux, opts.OwnersOptions)

	propertiesEndpoints.RegisterHandlers(mux, opts.PropsOptions)

	paymentEndpoints.RegisterHandlers(mux, opts.PayOptions)

	transactionsEndpoints.RegisterHandlers(mux, opts.TransOptions)

	accountsEndpoints.RegisterHandlers(mux, opts.AccountsOptions)

	authEndpoints.RegisterHandlers(mux, opts.AuthOptions)

	invoiceEndpoints.RegisterHandlers(mux, opts.InvoiceOptions)

	metricsEndpoints.RegisterHandlers(mux, opts.StatsOptions)
}
