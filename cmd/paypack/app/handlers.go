package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/auth"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/feedback"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/health"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/metrics"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/owners"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/payment"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/properties"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/scheduler"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/users"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/ussd"
	"github.com/nshimiyimanaamani/paypack-backend/api/http/version"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// HandlerOptions ...
type HandlerOptions struct {
	AccountsOptions  *accounts.HandlerOpts
	AuthOptions      *auth.HandlerOpts
	FeedbackOptions  *feedback.HandlerOpts
	NotifOptions     *notifs.HandlerOpts
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
		Logger:        lggr,
		Service:       services.Payment,
		Authenticator: services.Auth,
		Repository:    services.PaymentStore,
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
	notifOpts := &notifs.HandlerOpts{
		Logger:        lggr,
		Service:       services.Notifications,
		Authenticator: services.Auth,
	}
	ussdOpts := &ussd.HandlerOpts{
		Service: services.USSD,
		Logger:  lggr,
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

	notifs.RegisterHandlers(mux, opts.NotifOptions)

	ussd.RegisterHandlers(mux, opts.USSDOptions)

	scheduler.RegisterHandlers(mux, opts.SchedulerOptions)
}
