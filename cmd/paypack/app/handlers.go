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
	"github.com/rugwirobaker/paypack-backend/api/http/ussd"
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
	USSDOptions      *ussd.HandlerOpts
}

// NewHandlerOptions ...
func NewHandlerOptions(services *Services, lggr *log.Logger) *HandlerOptions {
	feedOpts := &feedback.HandlerOpts{
		Service: services.Feedback,
		Logger:  lggr,
	}
	proOpts := &properties.HandlerOpts{
		Service: services.Properties,
		Logger:  lggr,
	}

	ownersOpts := &owners.HandlerOpts{
		Service: services.Owners,
		Logger:  lggr,
	}
	paymentOpts := &payment.HandlerOpts{
		Service: services.Payment,
		Logger:  lggr,
	}

	transOpts := &transactions.HandlerOpts{
		Service: services.Transactions,
		Logger:  lggr,
	}

	usersOpts := &users.HandlerOpts{
		Service: services.Users,
		Logger:  lggr,
	}

	accountsOpts := &accounts.HandlerOpts{
		Service: services.Accounts,
		Logger:  lggr,
	}

	authOpts := &auth.HandlerOpts{
		Service: services.Auth,
		Logger:  lggr,
	}
	invOpts := &invoices.HandlerOpts{
		Service: services.Invoices,
		Logger:  lggr,
	}
	statsOpts := &metrics.HandlerOpts{
		Service: services.Stats,
		Logger:  lggr,
	}
	notifOpts := &notifications.HandlerOpts{
		Service: services.Notifications,
		Logger:  lggr,
	}
	ussdOpts := &ussd.HandlerOpts{
		Service: services.USSD,
		Logger:  lggr,
	}
	scOptions := &scheduler.HandlerOpts{
		Service: services.Scheduler,
		Logger:  lggr,
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
		USSDOptions:      ussdOpts,
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

	notifications.RegisterHandlers(mux, opts.NotifOptions)

	ussd.RegisterHandlers(mux, opts.USSDOptions)

	scheduler.RegisterHandlers(mux, opts.SchedulerOptions)
}
