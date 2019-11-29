package app

import (
	"net/http"

	"github.com/gorilla/mux"
	feedbackEndpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/api/http/health"
	ownersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/owners"
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
		PropsOptions:    proOpts,
		OwnersOptions:   ownersOpts,
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
}
