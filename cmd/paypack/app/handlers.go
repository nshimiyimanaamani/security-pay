package app

import (
	"net/http"

	"github.com/gorilla/mux"
	feedbackEndpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/api/http/health"
	ownersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/owners"
	paymentEndpoints "github.com/rugwirobaker/paypack-backend/api/http/payment"
	propertiesEndpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	transactionsEndpoints "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	usersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/api/http/version"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// HandlerOptions ...
type HandlerOptions struct {
	FeedbackOptions *feedbackEndpoints.HandlerOpts
	OwnersOptions   *ownersEndpoints.HandlerOpts
	PayOptions      *paymentEndpoints.HandlerOpts
	PropsOptions    *propertiesEndpoints.HandlerOpts
	TransOptions    *transactionsEndpoints.HandlerOpts
	UsersOptions    *usersEndpoints.HandlerOpts
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

	opts := &HandlerOptions{
		FeedbackOptions: feedOpts,
		OwnersOptions:   ownersOpts,
		PropsOptions:    proOpts,
		PayOptions:      paymentOpts,
		TransOptions:    transOpts,
		UsersOptions:    usersOpts,
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
}
