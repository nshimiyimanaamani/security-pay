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
	"github.com/rugwirobaker/paypack-backend/api/http/scheduler"
	"github.com/rugwirobaker/paypack-backend/api/http/transactions"
	"github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/api/http/version"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// HandlerOptions ...
type HandlerOptions struct {
	AccountsOptions  *accounts.HandlerOpts
	AuthOptions      *auth.HandlerOpts
	FeedbackOptions  *feedback.HandlerOpts
	NotifOptions     *notifications.HandlerOpts
	OwnersOptions    *owners.HandlerOpts
	PayOptions       *payment.HandlerOpts
	PropsOptions     *properties.HandlerOpts
	TransOptions     *transactions.HandlerOpts
	UsersOptions     *users.HandlerOpts
	InvoiceOptions   *invoices.HandlerOpts
	StatsOptions     *metrics.HandlerOpts
	SchedulerOptions *scheduler.HandlerOpts
}

// NewHandlerOptions ...
func NewHandlerOptions(services *Services, lggr *log.Logger) *HandlerOptions {
	feedOpts := &feedback.HandlerOpts{
		Logger:        lggr,
		Service:       services.Feedback,
		Authenticator: services.Auth,
	}
	proOpts := &properties.HandlerOpts{
		Logger:        lggr,
		Service:       services.Properties,
		Authenticator: services.Auth,
	}

	ownersOpts := &owners.HandlerOpts{
		Logger:        lggr,
		Service:       services.Owners,
		Authenticator: services.Auth,
	}
	paymentOpts := &payment.HandlerOpts{
		Service: services.Payment,
		Logger:  lggr,
	}

	transOpts := &transactions.HandlerOpts{
		Logger:        lggr,
		Service:       services.Transactions,
		Authenticator: services.Auth,
	}

	usersOpts := &users.HandlerOpts{
		Logger:        lggr,
		Service:       services.Users,
		Authenticator: services.Auth,
	}

	accountsOpts := &accounts.HandlerOpts{
		Logger:        lggr,
		Service:       services.Accounts,
		Authenticator: services.Auth,
	}

	authOpts := &auth.HandlerOpts{
		Service: services.Auth,
		Logger:  lggr,
	}
	invOpts := &invoices.HandlerOpts{
		Logger:        lggr,
		Service:       services.Invoices,
		Authenticator: services.Auth,
	}
	statsOpts := &metrics.HandlerOpts{
		Logger:        lggr,
		Service:       services.Stats,
		Authenticator: services.Auth,
	}
	notifOpts := &notifications.HandlerOpts{
		Logger:        lggr,
		Service:       services.Notifications,
		Authenticator: services.Auth,
	}
	scOptions := &scheduler.HandlerOpts{
		Logger:        lggr,
		Service:       services.Scheduler,
		Authenticator: services.Auth,
	}

	opts := &HandlerOptions{
		AuthOptions:      authOpts,
		AccountsOptions:  accountsOpts,
		FeedbackOptions:  feedOpts,
		OwnersOptions:    ownersOpts,
		PropsOptions:     proOpts,
		PayOptions:       paymentOpts,
		TransOptions:     transOpts,
		UsersOptions:     usersOpts,
		InvoiceOptions:   invOpts,
		StatsOptions:     statsOpts,
		NotifOptions:     notifOpts,
		SchedulerOptions: scOptions,
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
	mux.HandleFunc("/panic", health.Panic).Methods(http.MethodGet)

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

	notifications.RegisterHandlers(mux, opts.NotifOptions)

	scheduler.RegisterHandlers(mux, opts.SchedulerOptions)
}
